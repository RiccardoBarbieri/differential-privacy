package command

import (
	"github.com/spf13/cobra"
	"godp"
	"godp/runs"
)

var AllCmd = &cobra.Command{
	Use:     "all",
	Short:   "Run all the operations in the pipeline",
	RunE:    runs.RunAll,
	Args:    cobra.NoArgs,
	PreRunE: healthcaredp.InitEnvironment,
}

func init() {
}
