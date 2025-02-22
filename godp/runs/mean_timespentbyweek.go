package runs

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"healthcaredp"
	"healthcaredp/aggregations"
	"healthcaredp/model/utils"
	"strings"
)

var stayByWeekMeanOutputCsv string
var stayByWeekMeanOutputCsvDp string

func MeanStayByWeek(scope beam.Scope, outputCsv string, generateClear bool, admissions beam.PCollection) {
	scope = scope.Scope("MeanStayByWeek")

	baseOutputName := strings.TrimSuffix(outputCsv, ".csv")
	stayByWeekMeanOutputCsv = baseOutputName + "_mean_stay_by_week.csv"
	stayByWeekMeanOutputCsvDp = baseOutputName + "_mean_stay_by_week_dp.csv"

	if generateClear {
		meanStayByWeek := aggregations.MeanStayByWeek(scope, admissions)
		utils.WriteOutput(scope, meanStayByWeek, stayByWeekMeanOutputCsv)
	}

	meanStayByWeekDp := aggregations.MeanStayByWeekDp(scope, admissions, healthcaredp.Budget)
	utils.WriteOutput(scope, meanStayByWeekDp, stayByWeekMeanOutputCsvDp)
}

func MeanStayByWeekWriteHeaders(generateClear bool) {
	if generateClear {
		utils.WriteHeaders(stayByWeekMeanOutputCsv, "Week", "Mean Stay (Days)")
	}
	utils.WriteHeaders(stayByWeekMeanOutputCsvDp, "Week", "Mean Stay (Days) (DP)")
}
