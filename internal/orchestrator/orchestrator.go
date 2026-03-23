package orchestrator

import (
	"fmt"

	"github.com/tacomilkshake/orb-optimizer/internal/connector"
	"github.com/tacomilkshake/orb-optimizer/internal/store"
)

// Orchestrator runs automated sweep tests.
type Orchestrator struct {
	store     *store.Store
	connector connector.APConnector
}

// New creates a sweep orchestrator.
func New(s *store.Store, c connector.APConnector) *Orchestrator {
	return &Orchestrator{store: s, connector: c}
}

// Run executes the sweep matrix. Returns test IDs for each combination.
func (o *Orchestrator) Run(cfg SweepConfig) ([]int64, error) {
	var testIDs []int64

	for _, ch := range cfg.Channels {
		for _, w := range cfg.Widths {
			name := fmt.Sprintf("%s_ch%d_%dMHz", cfg.Band, ch, w)

			// Set the radio
			err := o.connector.SetRadio(cfg.APMAC, connector.RadioConfig{
				Channel:  ch,
				WidthMHz: w,
			})
			if err != nil {
				return testIDs, fmt.Errorf("set radio ch%d/%dMHz: %w", ch, w, err)
			}

			fmt.Printf("[sweep] Radio set to ch%d/%dMHz, settling for %s...\n", ch, w, cfg.Settle)
			// TODO: time.Sleep(cfg.Settle) — but also need to keep collecting during settle

			// Begin test
			id, err := o.store.BeginTest(store.BeginTestParams{
				Name:       name,
				Channel:    ch,
				WidthMHz:   w,
				FreqMHz:    0, // TODO: derive from channel
				APPlatform: o.connector.Name(),
			})
			if err != nil {
				return testIDs, fmt.Errorf("begin test %s: %w", name, err)
			}
			testIDs = append(testIDs, id)

			fmt.Printf("[sweep] Test #%d: %s — collecting for %s...\n", id, name, cfg.Window)
			// TODO: time.Sleep(cfg.Window) — collector must be running in parallel

			// End test
			_, _, err = o.store.EndTest(id)
			if err != nil {
				return testIDs, fmt.Errorf("end test %s: %w", name, err)
			}

			fmt.Printf("[sweep] Test #%d: %s — ended\n", id, name)
		}
	}

	return testIDs, nil
}
