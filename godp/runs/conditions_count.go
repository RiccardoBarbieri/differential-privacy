package runs

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"healthcaredp"
	"healthcaredp/aggregations"
	"healthcaredp/model/utils"
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
		utils.WriteOutput(scope, conditionsCount, conditionsCountOutputCsv)
	}

	conditionsCountDp := aggregations.CountConditionsDp(scope, admissions, healthcaredp.Budget)
	utils.WriteOutput(scope, conditionsCountDp, conditionsCountOutputCsvDp)
}

func ConditionsCountWriteHeaders(generateClear bool) {
	if generateClear {
		utils.WriteHeaders(conditionsCountOutputCsv, "Medical Condition", "Count")
	}
	utils.WriteHeaders(conditionsCountOutputCsvDp, "Medical Condition", "Count(DP)")
}
