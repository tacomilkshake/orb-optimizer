// Package orchestrator implements the sweep test matrix logic.
package orchestrator

import "time"

// SweepConfig defines parameters for an automated test sweep.
type SweepConfig struct {
	APMAC    string
	Band     string
	Channels []int
	Widths   []int
	Window   time.Duration // test duration per combination
	Settle   time.Duration // wait after radio change before testing
}
