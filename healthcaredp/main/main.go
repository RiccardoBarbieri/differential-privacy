package main

import (
	"context"
	"flag"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/runners/direct"
	log "github.com/golang/glog"
	"healthcaredp"
	"healthcaredp/aggregations"
	"healthcaredp/utils"
	"strings"
)

var (
	generateClear = flag.Bool("generate_clear", false, "Generate clear  along with conditions and test results")
	inputCsv      = flag.String("input_csv", "", "Name of the csv file that contains the healthcare data")
	outputCsv     = flag.String("output_csv", "", "Base name of the output csv file that will contain output data (ex. output.csv)")
	outputClean   = flag.String("output_clean", "", "Name of the output csv file that will contain the cleaned dataset")
)

func main() {
	flag.Parse()
	if *inputCsv == "" {
		log.Exit("Input csv file is required.")
	}
	if *outputCsv == "" {
		log.Exit("Output csv file is required.")
	}
	if *outputClean == "" {
		log.Exit("Output clean csv file is required.")
	}
	baseOutputName := strings.TrimSuffix(*outputCsv, ".csv")
	ccOutputCsv := baseOutputName + "_conditions_count.csv"
	ccOutputCsvDp := baseOutputName + "_conditions_count_dp.csv"
	ctrOutputCsv := baseOutputName + "_testresults_count.csv"
	ctrOutputCsvDp := baseOutputName + "_testresults_count_dp.csv"

	beam.Init()

	pipeline := beam.NewPipeline()
	scope := pipeline.Root()

	healtcaredp.InitBudgetSplits(1, 2)
	globalPrivacySpec := healtcaredp.GlobalPrivacySpec

	admissions := utils.ReadInput(scope, *inputCsv)
	admissionsCleaned := healtcaredp.CleanDataset(scope, admissions)
	utils.WriteOutput(scope, admissionsCleaned, *outputClean)

	if *generateClear {
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

	utils.WriteHeaders(*outputClean, utils.StructCsvHeaders(healtcaredp.Admission{})...)
	if *generateClear {
		utils.WriteHeaders(ccOutputCsv, "Medical Condition", "Count")
		utils.WriteHeaders(ctrOutputCsv, "Test Results", "Count")
	}
	utils.WriteHeaders(ccOutputCsvDp, "Medical Condition", "Count(DP)")
	utils.WriteHeaders(ctrOutputCsvDp, "Test Results", "Count(DP)")

}
