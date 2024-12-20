package runs

import (
	"context"
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/runners/direct"
	log "github.com/golang/glog"
	"github.com/spf13/cobra"
	healtcaredp "healthcaredp"
	"healthcaredp/aggregations"
	"healthcaredp/utils"
	"strings"
)

func RunAll(cmd *cobra.Command, args []string) error {
	inputCsv, _ := cmd.Parent().PersistentFlags().GetString("input-csv")
	outputCsv, _ := cmd.Parent().PersistentFlags().GetString("output-csv")
	outputClean, _ := cmd.Parent().PersistentFlags().GetString("output-clean")
	generateClear, _ := cmd.Parent().PersistentFlags().GetBool("generate-non-dp")

	if inputCsv == "" {
		return fmt.Errorf("input-csv flag is required")
	}
	if outputCsv == "" {
		return fmt.Errorf("output-csv flag is required")
	}
	if outputClean == "" {
		return fmt.Errorf("output-clean flag is required")
	}

	healtcaredp.Budget.InitAllBudgetShares(1, 2)

	beam.Init()
	RunSum(inputCsv, outputCsv, outputClean, generateClear)
	return nil
}

func RunSum(inputCsv string, outputCsv string, outputClean string, generateClear bool) {

	baseOutputName := strings.TrimSuffix(outputCsv, ".csv")
	ccOutputCsv := baseOutputName + "_conditions_count.csv"
	ccOutputCsvDp := baseOutputName + "_conditions_count_dp.csv"
	ctrOutputCsv := baseOutputName + "_testresults_count.csv"
	ctrOutputCsvDp := baseOutputName + "_testresults_count_dp.csv"

	pipeline := beam.NewPipeline()
	scope := pipeline.Root()

	globalPrivacySpec := healtcaredp.Budget.PrivacySpec

	admissions := utils.ReadInput(scope, inputCsv)
	admissionsCleaned := healtcaredp.CleanDataset(scope, admissions)
	utils.WriteOutput(scope, admissionsCleaned, outputClean)

	if generateClear {
		conditionsCount := aggregations.CountConditions(scope, admissionsCleaned)
		testResultsCount := aggregations.CountTestResults(scope, admissionsCleaned)
		utils.WriteOutput(scope, conditionsCount, ccOutputCsv)
		utils.WriteOutput(scope, testResultsCount, ctrOutputCsv)
	}

	conditionsCountDp := aggregations.CountConditionsDp(scope, admissionsCleaned, globalPrivacySpec, healtcaredp.Budget.GetBudgetShare("CountConditionsDp"))
	testResultsCountDp := aggregations.CountTestResultsDp(scope, admissionsCleaned, globalPrivacySpec, healtcaredp.Budget.GetBudgetShare("CountTestResultsDp"))
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
