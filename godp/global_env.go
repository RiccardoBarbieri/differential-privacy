package model

import (
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	log "github.com/golang/glog"
	"github.com/spf13/cobra"
)

type IOArgs struct {
	InputCsv      string
	OutputCsv     string
	OutputClean   string
	GenerateNonDp bool
	Args          []string
}

// CurrentIOArgs stores the flags, parsed by PreRun on RootCmd
var CurrentIOArgs IOArgs = IOArgs{}

// GlobalScope is used to store a global beam.Scope, init by PreRun RootCmd
var GlobalScope beam.Scope

// GlobalPipeline is used to store a global beam.Pipeline, init by PreRun RootCmd
var GlobalPipeline *beam.Pipeline

var AdmissionsCleaned beam.PCollection

func InitEnvironment(cmd *cobra.Command, args []string) (err error) {
	log.Info("Initializing environment")
	var _ error
	CurrentIOArgs.InputCsv, _ = cmd.Parent().PersistentFlags().GetString("input-csv")
	CurrentIOArgs.OutputCsv, _ = cmd.Parent().PersistentFlags().GetString("output-csv")
	CurrentIOArgs.OutputClean, _ = cmd.Parent().PersistentFlags().GetString("output-clean")
	CurrentIOArgs.GenerateNonDp, _ = cmd.Parent().PersistentFlags().GetBool("generate-non-dp")
	log.Infof("input-csv: %s, output-csv: %s, output-clean: %s, generate-non-dp: %t\n", CurrentIOArgs.InputCsv, CurrentIOArgs.OutputCsv, CurrentIOArgs.OutputClean, CurrentIOArgs.GenerateNonDp)

	if CurrentIOArgs.InputCsv == "" {
		log.Error("Input CSV file is empty")
		return fmt.Errorf("input CSV file is required")
	}
	if CurrentIOArgs.OutputCsv == "" {
		log.Error("Output CSV file is empty")
		return fmt.Errorf("output CSV file is required")
	}
	if CurrentIOArgs.OutputClean == "" {
		log.Error("Output clean CSV file is empty")
		return fmt.Errorf("output clean CSV file is required")
	}

	beam.Init()
	GlobalPipeline = beam.NewPipeline()
	GlobalScope = GlobalPipeline.Root()

	AdmissionsCleaned = LoadCleanDataset(GlobalScope, CurrentIOArgs.InputCsv)
	WriteOutput(GlobalScope, AdmissionsCleaned, CurrentIOArgs.OutputClean)

	return nil
}
