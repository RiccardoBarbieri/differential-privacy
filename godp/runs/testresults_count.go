package runs

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"godp"
	"godp/aggregations"
	"godp/model/utils"
	"strings"
)

var testresultsCountOutputCsv string
var testresultsCountOutputCsvDp string

func RunTestResultsCount(scope beam.Scope, outputCsv string, generateClear bool, admissions beam.PCollection) {
	scope = scope.Scope("RunTestResultsCount")

	baseOutputName := strings.TrimSuffix(outputCsv, ".csv")
	testresultsCountOutputCsv = baseOutputName + "_testresults_count.csv"
	testresultsCountOutputCsvDp = baseOutputName + "_testresults_count_dp.csv"

	if generateClear {
		testResultsCount := aggregations.CountTestResults(scope, admissions)
		utils.WriteOutput(scope, testResultsCount, testresultsCountOutputCsv)
	}

	testResultsCountDp := aggregations.CountTestResultsDp(scope, admissions, healthcaredp.Budget)
	utils.WriteOutput(scope, testResultsCountDp, testresultsCountOutputCsvDp)
}

func TestResultsCountWriteHeaders(generateClear bool) {
	if generateClear {
		utils.WriteHeaders(testresultsCountOutputCsv, "Test Results", "Count")
	}
	utils.WriteHeaders(testresultsCountOutputCsvDp, "Test Results", "Count(DP)")
}
