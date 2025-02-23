package command

import (
	"github.com/spf13/cobra"
	"godp/runs"
)

var FromFileCmd = &cobra.Command{
	Use:   "fromfile",
	Short: "Read input from a YAML file",
	RunE:  runs.RunFromFile,
}

func init() {

	FromFileCmd.PersistentFlags().String("file", "", "Name of the yaml file to read the configuration from")

	var err error
	err = FromFileCmd.MarkPersistentFlagRequired("file")
	if err != nil {
		return
	}
	err = FromFileCmd.MarkPersistentFlagFilename("file", "yaml", "yml")
	if err != nil {
		return
	}
}
