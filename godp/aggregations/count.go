package aggregations

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/transforms/stats"
	"github.com/google/differential-privacy/privacy-on-beam/v3/pbeam"
	"godp"
	"godp/model"
)

func init() {
	register.Function1x1[model.Admission, string](extractMedicalCondition)
	register.Function1x1[model.Admission, string](extractTestResult)
}

func extractMedicalCondition(admission model.Admission) string {
	return admission.MedicalCondition
}

func extractTestResult(admission model.Admission) string {
	return admission.TestResults
}

func CountConditions(scope beam.Scope, col beam.PCollection) beam.PCollection {
	scope = scope.Scope("CountConditions")
	conditions := beam.ParDo(scope, extractMedicalCondition, col)
	conditionsCount := stats.Count(scope, conditions)
	return conditionsCount
}

func CountConditionsDp(scope beam.Scope, col beam.PCollection, budget healthcaredp.DpBudget) beam.PCollection {
	operation := "CountConditionsDp"
	scope = scope.Scope(operation)
	pCol := pbeam.MakePrivateFromStruct(scope, col, budget.PrivacySpec, "Name")

	pConditions := pbeam.ParDo(scope, extractMedicalCondition, pCol)
	pConditionsCount := pbeam.Count(scope, pConditions, pbeam.CountParams{
		PartitionSelectionParams: pbeam.PartitionSelectionParams{
			Epsilon: budget.GetBudgetShare(operation).PartitionEpsilon,
			Delta:   budget.GetBudgetShare(operation).PartitionDelta,
		},
		AggregationEpsilon:       budget.GetBudgetShare(operation).AggregationEpsilon,
		MaxPartitionsContributed: 6,
		MaxValue:                 24,
	})
	return pConditionsCount
}

func CountTestResults(scope beam.Scope, col beam.PCollection) beam.PCollection {
	scope = scope.Scope("CountTestResults")
	conditions := beam.ParDo(scope, extractTestResult, col)
	conditionsCount := stats.Count(scope, conditions)
	return conditionsCount
}

func CountTestResultsDp(scope beam.Scope, col beam.PCollection, budget healthcaredp.DpBudget) beam.PCollection {
	operation := "CountTestResultsDp"
	scope = scope.Scope(operation)
	pCol := pbeam.MakePrivateFromStruct(scope, col, budget.PrivacySpec, "Name")

	pTestResults := pbeam.ParDo(scope, extractTestResult, pCol)
	pTestResultsCount := pbeam.Count(scope, pTestResults, pbeam.CountParams{
		PartitionSelectionParams: pbeam.PartitionSelectionParams{
			Epsilon: budget.GetBudgetShare(operation).PartitionEpsilon,
			Delta:   budget.GetBudgetShare(operation).PartitionDelta,
		},
		AggregationEpsilon:       budget.GetBudgetShare(operation).AggregationEpsilon,
		MaxPartitionsContributed: 3,
		MaxValue:                 24,
	})
	return pTestResultsCount
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
