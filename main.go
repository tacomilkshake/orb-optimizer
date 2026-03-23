// Package main provides the orb-optimizer CLI.
package main

import (
	"os"

	"github.com/tacomilkshake/orb-optimizer/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
