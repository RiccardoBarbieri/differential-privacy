package command

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "healthcaredp",
	Short:   "A pipeline to anonymize healthcare data using differential privacy techniques with privacy-on-beam",
	Args:    cobra.NoArgs,
	Example: "healthcaredp all --input-csv input.csv --output-csv output.csv --output-clean output_clean.csv",
}

func init() {

	RootCmd.AddCommand(AllCmd)
	RootCmd.AddCommand(CountCmd)
	RootCmd.AddCommand(MeanCmd)

	RootCmd.PersistentFlags().String("input-csv", "", "Name of the csv file that contains the healthcare data")
	RootCmd.PersistentFlags().String("output-csv", "", "Base name of the output csv file that will contain output data (ex. output.csv)")
	RootCmd.PersistentFlags().String("output-clean", "", "Name of the output csv file that will contain the cleaned dataset")
	RootCmd.PersistentFlags().Bool("generate-non-dp", false, "Generate non-differentially private data")

	var err error
	err = RootCmd.MarkFlagRequired("input-csv")
	if err != nil {
		return
	}
	err = RootCmd.MarkFlagRequired("output-csv")
	if err != nil {
		return
	}
	err = RootCmd.MarkFlagRequired("output-clean")
	if err != nil {
		return
	}
	RootCmd.SilenceUsage = true
}
