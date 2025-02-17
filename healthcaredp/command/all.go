package command

import (
	"github.com/spf13/cobra"
	"healthcaredp/model"
	"healthcaredp/runs"
)

var AllCmd = &cobra.Command{
	Use:     "all",
	Short:   "Run all the operations in the pipeline",
	RunE:    runs.RunAll,
	Args:    cobra.NoArgs,
	PreRunE: model.InitEnvironment,
}

func init() {
}
