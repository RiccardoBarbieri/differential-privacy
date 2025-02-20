package aggregations

import (
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/transforms/stats"
	log "github.com/golang/glog"
	"github.com/google/differential-privacy/privacy-on-beam/v3/pbeam"
	"healthcaredp"
	"healthcaredp/model"
	"math"
	"strconv"
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

func MeanColumnByKey(scope beam.Scope, col beam.PCollection, op model.OperationType, bd healthcaredp.DpBudget) (*beam.PCollection, error) {
	scope = scope.Scope(op.OperationName)
	pCol := pbeam.MakePrivateFromStruct(scope, col, bd.PrivacySpec, "Id")
	if _, ok := model.TypesMap[op.Column]; !ok {
		return nil, fmt.Errorf("column type not specified for column: %s", op.Column)
	}
	if model.TypesMap[op.Column] != "float" && model.TypesMap[op.Column] != "int" {
		return nil, fmt.Errorf("unsupported column type: %s for %s operation", model.TypesMap[op.Column], op.OperationType)
	}

	pColumnValuesByKey := pbeam.ParDo(scope, func(struc model.ValuesStruct) (string, float64) {
		castedCol, err := strconv.ParseFloat(struc.Values[op.Column], 64)
		if err != nil {
			log.Fatalf("Failed to convert column value %s to type %s: %v", op.Column, model.TypesMap[op.Column], err)
		}
		return struc.Values[*op.KeyColumn], castedCol
	}, pCol)

	pValuesMeanByKey := pbeam.MeanPerKey(scope, pColumnValuesByKey, pbeam.MeanParams{
		MaxPartitionsContributed:     *op.PrivacyParams.MaxCategoriesContributed,
		MaxContributionsPerPartition: *op.PrivacyParams.MaxContributionsPerCategory,
		MinValue:                     *op.PrivacyParams.MinValue,
		MaxValue:                     *op.PrivacyParams.MaxValue,

		AggregationEpsilon: bd.GetBudgetShare(op.OperationName).AggregationEpsilon,
		AggregationDelta:   bd.GetBudgetShare(op.OperationName).AggregationDelta,
		PartitionSelectionParams: pbeam.PartitionSelectionParams{
			Epsilon: bd.GetBudgetShare(op.OperationName).PartitionEpsilon,
			Delta:   bd.GetBudgetShare(op.OperationName).PartitionDelta,
		},
	})
	return &pValuesMeanByKey, nil

}
