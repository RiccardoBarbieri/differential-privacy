package command

import (
	"github.com/spf13/cobra"
	"godp"
	"godp/runs"
)

var MeanCmd = &cobra.Command{
	Use:               "mean",
	Short:             "Run the mean operations in the pipeline",
	RunE:              runs.RunMeans,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: checkMeanArg,
	PreRunE:           healthcaredp.InitEnvironment,
}

func init() {

}

func checkMeanArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return healthcaredp.MeanOperations, cobra.ShellCompDirectiveNoFileComp
}
