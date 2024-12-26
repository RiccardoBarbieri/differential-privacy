package aggregations

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/transforms/stats"
	"github.com/google/differential-privacy/privacy-on-beam/v3/pbeam"
	"healthcaredp"
	"healthcaredp/model"
	"math"
)

func init() {
	register.Function1x2[model.Admission, int, float64](extractWeekStayDuration)
}

func extractWeekStayDuration(admission model.Admission) (int, float64) {
	_, week := admission.DateOfAdmission.ISOWeek()
	stayDurationDays := math.Round(admission.DischargeDate.Sub(admission.DateOfAdmission).Hours() / 24)
	return week, stayDurationDays
}

func MeanStayByWeek(scope beam.Scope, col beam.PCollection) beam.PCollection {
	scope = scope.Scope("MeanStayByWeek")
	stayDurationByWeek := beam.ParDo(scope, extractWeekStayDuration, col)
	meanStayDuration := stats.MeanPerKey(scope, stayDurationByWeek)
	return meanStayDuration
}

func MeanStayByWeekDp(scope beam.Scope, col beam.PCollection, budget healthcaredp.DpBudget) beam.PCollection {
	operation := "MeanStayByWeekDp"
	scope = scope.Scope(operation)
	pCol := pbeam.MakePrivateFromStruct(scope, col, budget.PrivacySpec, "Name")

	pStayDurationByWeek := pbeam.ParDo(scope, extractWeekStayDuration, pCol)
	pGroupedStayDuration := pbeam.MeanPerKey(scope, pStayDurationByWeek, pbeam.MeanParams{
		MaxPartitionsContributed:     21,
		MaxContributionsPerPartition: 1039,
		MinValue:                     1,
		MaxValue:                     30,

		AggregationEpsilon: budget.GetBudgetShare(operation).AggregationEpsilon,
		AggregationDelta:   budget.GetBudgetShare(operation).AggregationDelta,
	})
	return pGroupedStayDuration
}
