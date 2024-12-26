package cmds

import (
	"github.com/spf13/cobra"
	"healthcaredp"
	"healthcaredp/runs"
)

var AllCmd = &cobra.Command{
	Use:    "all",
	Short:  "Run all the operations in the pipeline",
	RunE:   runs.RunAll,
	Args:   cobra.NoArgs,
	PreRun: healthcaredp.InitEnvironment,
}

func init() {
}
