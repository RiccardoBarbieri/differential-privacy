package runs

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"healthcaredp"
	"healthcaredp/aggregations"
	"healthcaredp/utils"
	"strings"
)

var stayByWeekAvgOutputCsv string
var stayByWeekAvgOutputCsvDp string

func RunAvgStayByWeek(scope beam.Scope, outputCsv string, generateClear bool, admissions beam.PCollection) {
	scope = scope.Scope("RunAvgStayByWeek")

	baseOutputName := strings.TrimSuffix(outputCsv, ".csv")
	stayByWeekAvgOutputCsv = baseOutputName + "_avg_stay_by_week.csv"
	stayByWeekAvgOutputCsvDp = baseOutputName + "_avg_stay_by_week_dp.csv"

	if generateClear {
		meanStayByWeek := aggregations.MeanStayByWeek(scope, admissions)
		utils.WriteOutput(scope, meanStayByWeek, stayByWeekAvgOutputCsv)
	}

	meanStayByWeekDp := aggregations.MeanStayByWeekDp(scope, admissions, healthcaredp.Budget)
	utils.WriteOutput(scope, meanStayByWeekDp, stayByWeekAvgOutputCsvDp)
}

func AvgStayByWeekWriteHeaders(generateClear bool) {
	if generateClear {
		utils.WriteHeaders(stayByWeekAvgOutputCsv, "Week", "Average Stay (Days)")
	}
	utils.WriteHeaders(stayByWeekAvgOutputCsvDp, "Week", "Average Stay (Days) (DP)")
}
