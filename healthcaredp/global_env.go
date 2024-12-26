package healthcaredp

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	log "github.com/golang/glog"
	"github.com/spf13/cobra"
	"healthcaredp/utils"
)

type IOArgs struct {
	InputCsv      string
	OutputCsv     string
	OutputClean   string
	GenerateNonDp bool
}

// CurrentIOArgs stores the flags, parsed by PreRun on RootCmd
var CurrentIOArgs IOArgs = IOArgs{}

// GlobalScope is used to store a global beam.Scope, init by PreRun RootCmd
var GlobalScope beam.Scope

// GlobalPipeline is used to store a global beam.Pipeline, init by PreRun RootCmd
var GlobalPipeline *beam.Pipeline

var AdmissionsCleaned beam.PCollection

func InitEnvironment(cmd *cobra.Command, args []string) {
	log.Info("Initializing environment")
	CurrentIOArgs.InputCsv, _ = cmd.Parent().PersistentFlags().GetString("input-csv")
	CurrentIOArgs.OutputCsv, _ = cmd.Parent().PersistentFlags().GetString("output-csv")
	CurrentIOArgs.OutputClean, _ = cmd.Parent().PersistentFlags().GetString("output-clean")
	CurrentIOArgs.GenerateNonDp, _ = cmd.Parent().PersistentFlags().GetBool("generate-non-dp")
	log.Infof("input-csv: %s, output-csv: %s, output-clean: %s, generate-non-dp: %t\n", CurrentIOArgs.InputCsv, CurrentIOArgs.OutputCsv, CurrentIOArgs.OutputClean, CurrentIOArgs.GenerateNonDp)

	if CurrentIOArgs.InputCsv == "" {
		log.Fatal("Input CSV file is required")
	}
	if CurrentIOArgs.OutputCsv == "" {
		log.Fatal("Output CSV file is required")
	}
	if CurrentIOArgs.OutputClean == "" {
		log.Fatal("Output clean CSV file is required")
	}

	beam.Init()
	GlobalPipeline = beam.NewPipeline()
	GlobalScope = GlobalPipeline.Root()

	AdmissionsCleaned = utils.LoadCleanDataset(GlobalScope, CurrentIOArgs.InputCsv)
	utils.WriteOutput(GlobalScope, AdmissionsCleaned, CurrentIOArgs.OutputClean)
}
