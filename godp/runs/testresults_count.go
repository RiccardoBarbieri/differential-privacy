package runs

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"healthcaredp"
	"healthcaredp/aggregations"
	"healthcaredp/model"
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
		model.WriteOutput(scope, testResultsCount, testresultsCountOutputCsv)
	}

	testResultsCountDp := aggregations.CountTestResultsDp(scope, admissions, healthcaredp.Budget)
	model.WriteOutput(scope, testResultsCountDp, testresultsCountOutputCsvDp)
}

func TestResultsCountWriteHeaders(generateClear bool) {
	if generateClear {
		model.WriteHeaders(testresultsCountOutputCsv, "Test Results", "Count")
	}
	model.WriteHeaders(testresultsCountOutputCsvDp, "Test Results", "Count(DP)")
}
