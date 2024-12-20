package main

import (
	"context"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/runners/direct"
	log "github.com/golang/glog"
	"github.com/urfave/cli/v3"
	"healthcaredp"
	"healthcaredp/aggregations"
	"healthcaredp/utils"
	"os"
	"strings"
)

func main() {
	cmd := &cli.Command{
		Version:               "0.1.0",
		EnableShellCompletion: true,
		Name:                  "healthcaredp",
		Usage:                 "A pipeline to anonymize healthcare data with differential privacy",
		Commands: []*cli.Command{
			{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Run all the operations in the pipeline",
				Action:  runAll,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input-csv",
						Value:    "",
						Required: true,
						Usage:    "Name of the csv `file` that contains the healthcare data",
					},
					&cli.StringFlag{
						Name:     "output-csv",
						Value:    "",
						Required: true,
						Usage:    "Base name of the output csv `file` that will contain output data (ex. output.csv)",
					},
					&cli.StringFlag{
						Name:     "output-clean",
						Value:    "",
						Required: true,
						Usage:    "Name of the output csv `file` that will contain the cleaned dataset",
					},
					&cli.BoolFlag{
						Name:  "generate-non-dp",
						Value: false,
						Usage: "Generate non dp results along with conditions and test results (dev only)",
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Error running command: %v", err)
	}

}

func runAll(ctx context.Context, cmd *cli.Command) error {
	inputCsv := cmd.String("input-csv")
	outputCsv := cmd.String("output-csv")
	outputClean := cmd.String("output-clean")
	generateClear := cmd.Bool("generate-non-dp")
	beam.Init()
	runSum(inputCsv, outputCsv, outputClean, generateClear)
	return nil
}

func runSum(inputCsv string, outputCsv string, outputClean string, generateClear bool) {

	baseOutputName := strings.TrimSuffix(outputCsv, ".csv")
	ccOutputCsv := baseOutputName + "_conditions_count.csv"
	ccOutputCsvDp := baseOutputName + "_conditions_count_dp.csv"
	ctrOutputCsv := baseOutputName + "_testresults_count.csv"
	ctrOutputCsvDp := baseOutputName + "_testresults_count_dp.csv"

	pipeline := beam.NewPipeline()
	scope := pipeline.Root()

	healtcaredp.InitBudgetSplits(1, 2)
	globalPrivacySpec := healtcaredp.GlobalPrivacySpec

	admissions := utils.ReadInput(scope, inputCsv)
	admissionsCleaned := healtcaredp.CleanDataset(scope, admissions)
	utils.WriteOutput(scope, admissionsCleaned, outputClean)

	if generateClear {
		conditionsCount := aggregations.CountConditions(scope, admissionsCleaned)
		testResultsCount := aggregations.CountTestResults(scope, admissionsCleaned)
		utils.WriteOutput(scope, conditionsCount, ccOutputCsv)
		utils.WriteOutput(scope, testResultsCount, ctrOutputCsv)
	}

	conditionsCountDp := aggregations.CountConditionsDp(scope, admissionsCleaned, globalPrivacySpec, healtcaredp.GetBudgetShare("CountConditionsDp"))
	testResultsCountDp := aggregations.CountTestResultsDp(scope, admissionsCleaned, globalPrivacySpec, healtcaredp.GetBudgetShare("CountTestResultsDp"))
	utils.WriteOutput(scope, conditionsCountDp, ccOutputCsvDp)
	utils.WriteOutput(scope, testResultsCountDp, ctrOutputCsvDp)

	// Execute pipeline.
	_, err := direct.Execute(context.Background(), pipeline)
	if err != nil {
		log.Exitf("Execution of pipeline failed: %v", err)
	}

	utils.WriteHeaders(outputClean, utils.StructCsvHeaders(healtcaredp.Admission{})...)
	if generateClear {
		utils.WriteHeaders(ccOutputCsv, "Medical Condition", "Count")
		utils.WriteHeaders(ctrOutputCsv, "Test Results", "Count")
	}
	utils.WriteHeaders(ccOutputCsvDp, "Medical Condition", "Count(DP)")
	utils.WriteHeaders(ctrOutputCsvDp, "Test Results", "Count(DP)")
}
