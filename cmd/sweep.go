package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tacomilkshake/orb-optimizer/internal/orchestrator"
	"github.com/tacomilkshake/orb-optimizer/internal/report"
)

var sweepFlags struct {
	apMAC    string
	band     string
	channels string
	widths   string
	window   int
	settle   int
}

func newSweepCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sweep",
		Short: "Automated test matrix across channels and widths",
		Long: `Run automated tests across a matrix of channels and widths.
Requires SetRadio support from the AP connector.

Example:
  orb-optimizer sweep --ap-mac A8-29-48-E0-85-50 --band 6g \
    --channels 53,149,165,37 --widths 80,160,320 \
    --window 600 --settle 30`,
		RunE: runSweep,
	}

	cmd.Flags().StringVar(&sweepFlags.apMAC, "ap-mac", "", "AP MAC address (required)")
	cmd.Flags().StringVar(&sweepFlags.band, "band", "6g", "Band (2g, 5g, 6g)")
	cmd.Flags().StringVar(&sweepFlags.channels, "channels", "", "Comma-separated channel list (required)")
	cmd.Flags().StringVar(&sweepFlags.widths, "widths", "80,160", "Comma-separated width list in MHz")
	cmd.Flags().IntVar(&sweepFlags.window, "window", 600, "Test duration in seconds")
	cmd.Flags().IntVar(&sweepFlags.settle, "settle", 30, "Settle time after radio change in seconds")

	_ = cmd.MarkFlagRequired("ap-mac")
	_ = cmd.MarkFlagRequired("channels")

	return cmd
}

func runSweep(cmd *cobra.Command, args []string) error {
	if apConn == nil {
		return fmt.Errorf("sweep requires an AP connector (--ap-connector)")
	}

	channels, err := parseIntList(sweepFlags.channels)
	if err != nil {
		return fmt.Errorf("parse channels: %w", err)
	}
	widths, err := parseIntList(sweepFlags.widths)
	if err != nil {
		return fmt.Errorf("parse widths: %w", err)
	}

	cfg := orchestrator.SweepConfig{
		APMAC:    sweepFlags.apMAC,
		Band:     sweepFlags.band,
		Channels: channels,
		Widths:   widths,
		Window:   time.Duration(sweepFlags.window) * time.Second,
		Settle:   time.Duration(sweepFlags.settle) * time.Second,
	}

	fmt.Printf("[sweep] Matrix: %d channels x %d widths = %d tests\n",
		len(channels), len(widths), len(channels)*len(widths))
	fmt.Printf("[sweep] Window: %ds, Settle: %ds\n", sweepFlags.window, sweepFlags.settle)

	orch := orchestrator.New(db, apConn)
	testIDs, err := orch.Run(cfg)
	if err != nil {
		return fmt.Errorf("sweep: %w", err)
	}

	fmt.Printf("[sweep] Completed %d tests\n", len(testIDs))

	// Generate report
	rows, err := db.GetReportRows()
	if err != nil {
		return fmt.Errorf("generate report: %w", err)
	}
	report.PrintReport(os.Stdout, rows)

	return nil
}

func parseIntList(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", p, err)
		}
		result = append(result, n)
	}
	return result, nil
}
