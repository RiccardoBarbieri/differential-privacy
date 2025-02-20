package runs

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"healthcaredp"
	"healthcaredp/aggregations"
	"healthcaredp/model"
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
		model.WriteOutput(scope, meanStayByWeek, stayByWeekMeanOutputCsv)
	}

	meanStayByWeekDp := aggregations.MeanStayByWeekDp(scope, admissions, healthcaredp.Budget)
	model.WriteOutput(scope, meanStayByWeekDp, stayByWeekMeanOutputCsvDp)
}

func MeanStayByWeekWriteHeaders(generateClear bool) {
	if generateClear {
		model.WriteHeaders(stayByWeekMeanOutputCsv, "Week", "Mean Stay (Days)")
	}
	model.WriteHeaders(stayByWeekMeanOutputCsvDp, "Week", "Mean Stay (Days) (DP)")
}
