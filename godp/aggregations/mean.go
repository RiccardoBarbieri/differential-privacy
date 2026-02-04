package aggregations

import (
	"fmt"
	"godp"
	"godp/model"
	"strconv"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	log "github.com/golang/glog"
	"github.com/google/differential-privacy/privacy-on-beam/v3/pbeam"
)

func init() {
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
