package runs

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"healthcaredp"
	"healthcaredp/aggregations"
	"healthcaredp/model"
	"strings"
)

var conditionsCountOutputCsv string
var conditionsCountOutputCsvDp string

func RunConditionsCount(scope beam.Scope, outputCsv string, generateClear bool, admissions beam.PCollection) {
	scope = scope.Scope("RunConditionsCount")

	baseOutputName := strings.TrimSuffix(outputCsv, ".csv")
	conditionsCountOutputCsv = baseOutputName + "_conditions_count.csv"
	conditionsCountOutputCsvDp = baseOutputName + "_conditions_count_dp.csv"

	if generateClear {
		conditionsCount := aggregations.CountConditions(scope, admissions)
		model.WriteOutput(scope, conditionsCount, conditionsCountOutputCsv)
	}

	conditionsCountDp := aggregations.CountConditionsDp(scope, admissions, healthcaredp.Budget)
	model.WriteOutput(scope, conditionsCountDp, conditionsCountOutputCsvDp)
}

func ConditionsCountWriteHeaders(generateClear bool) {
	if generateClear {
		model.WriteHeaders(conditionsCountOutputCsv, "Medical Condition", "Count")
	}
	model.WriteHeaders(conditionsCountOutputCsvDp, "Medical Condition", "Count(DP)")
}
