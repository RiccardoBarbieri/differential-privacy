package command

import (
	"github.com/spf13/cobra"
	"healthcaredp"
	"healthcaredp/runs"
)

var MeanCmd = &cobra.Command{
	Use:               "mean",
	Short:             "Run the mean operations in the pipeline",
	RunE:              runs.RunMean,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: checkMeanArg,
	PreRun:            healthcaredp.InitEnvironment,
}

func init() {

}

func checkMeanArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return healthcaredp.MeanOperations, cobra.ShellCompDirectiveNoFileComp
}
