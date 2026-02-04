package aggregations

import (
	"godp"
	"godp/model"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/google/differential-privacy/privacy-on-beam/v3/pbeam"
)

func init() {
}

func CountColumn(scope beam.Scope, col beam.PCollection, op model.OperationType, bd healthcaredp.DpBudget) (*beam.PCollection, error) {
	scope = scope.Scope(op.OperationName)
	pCol := pbeam.MakePrivateFromStruct(scope, col, bd.PrivacySpec, "Id")

	pColumnValues := pbeam.ParDo(scope,
		func(struc model.ValuesStruct) string {
			return struc.Values[op.Column]
		}, pCol)
	pColumnValuesCount := pbeam.Count(scope, pColumnValues,
		pbeam.CountParams{
			PartitionSelectionParams: pbeam.PartitionSelectionParams{
				Epsilon: bd.GetBudgetShare(op.OperationName).PartitionEpsilon,
				Delta:   bd.GetBudgetShare(op.OperationName).PartitionDelta,
			},
			AggregationEpsilon:       bd.GetBudgetShare(op.OperationName).AggregationEpsilon,
			AggregationDelta:         bd.GetBudgetShare(op.OperationName).AggregationDelta,
			MaxPartitionsContributed: *op.PrivacyParams.MaxCategoriesContributed,
			MaxValue:                 *op.PrivacyParams.MaxContributions,
		})
	return &pColumnValuesCount, nil
}
