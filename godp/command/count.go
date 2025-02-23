package command

import (
	"github.com/spf13/cobra"
	"godp"
	"godp/runs"
)

var CountCmd = &cobra.Command{
	Use:               "count",
	Short:             "Run the count operations in the pipeline",
	RunE:              runs.RunCounts,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: checkCountArg,
	PreRunE:           healthcaredp.InitEnvironment,
}

func init() {

}

func checkCountArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return healthcaredp.CountOperations, cobra.ShellCompDirectiveNoFileComp
}
