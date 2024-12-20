package cmd

import (
	"github.com/spf13/cobra"
	"healthcaredp/runs"
)

var AllCmd = &cobra.Command{
	Use:   "all",
	Short: "Run all the operations in the pipeline",
	RunE:  runs.RunAll,
}

func init() {
}
