package cmds

import (
	"github.com/spf13/cobra"
	"healthcaredp"
	"healthcaredp/runs"
)

var AvgCmd = &cobra.Command{
	Use:               "avg",
	Short:             "Run the avg operations in the pipeline",
	RunE:              runs.RunAvg,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: checkAvgArg,
	PreRun:            healthcaredp.InitEnvironment,
}

func init() {

}

func checkAvgArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return healthcaredp.AvgOperations, cobra.ShellCompDirectiveNoFileComp
}
