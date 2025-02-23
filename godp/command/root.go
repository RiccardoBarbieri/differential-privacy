package command

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Version: "1.0.0",
	Use:     "godp",
	Short:   "A pipeline to anonymize csv datasets using differential privacy techniques with privacy-on-beam",
	Args:    cobra.NoArgs,
	Example: "godp all --input-csv input.csv --output-csv output.csv --output-clean output_clean.csv",
}

func init() {

	RootCmd.AddCommand(AllCmd)
	RootCmd.AddCommand(CountCmd)
	RootCmd.AddCommand(MeanCmd)
	RootCmd.AddCommand(FromFileCmd)

	RootCmd.PersistentFlags().String("input-csv", "", "Name of the csv file that contains the csv dataset")
	RootCmd.PersistentFlags().String("output-csv", "", "Base name of the output csv file that will contain output data (ex. output.csv)")
	RootCmd.PersistentFlags().String("output-clean", "", "Name of the output csv file that will contain the cleaned dataset")
	RootCmd.PersistentFlags().Bool("generate-non-dp", false, "Generate non-differentially private data")

	//var err error
	//err = RootCmd.MarkPersistentFlagRequired("input-csv")
	//if err != nil {
	//	return
	//}
	//err = RootCmd.MarkPersistentFlagRequired("output-csv")
	//if err != nil {
	//	return
	//}
	//err = RootCmd.MarkPersistentFlagRequired("output-clean")
	//if err != nil {
	//	return
	//}
	RootCmd.SilenceUsage = true
}
