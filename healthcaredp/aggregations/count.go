package aggregations

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/transforms/stats"
	"github.com/google/differential-privacy/privacy-on-beam/v3/pbeam"
	"healthcaredp"
)

func init() {
	register.Function1x1[healtcaredp.Admission, string](extractMedicalCondition)
	register.Function1x1[healtcaredp.Admission, string](extractTestResult)
}

func extractMedicalCondition(admission healtcaredp.Admission) string {
	return admission.MedicalCondition
}

func extractTestResult(admission healtcaredp.Admission) string {
	return admission.TestResults
}

func CountConditions(scope beam.Scope, col beam.PCollection) beam.PCollection {
	scope = scope.Scope("CountConditions")
	conditions := beam.ParDo(scope, extractMedicalCondition, col)
	conditionsCount := stats.Count(scope, conditions)
	return conditionsCount
}

func CountConditionsDp(scope beam.Scope, col beam.PCollection, pSpec *pbeam.PrivacySpec, budgetShare healtcaredp.DpBudgetShare) beam.PCollection {
	scope = scope.Scope("PrivateCountConditions")
	pCol := pbeam.MakePrivateFromStruct(scope, col, pSpec, "Name")

	pConditions := pbeam.ParDo(scope, extractMedicalCondition, pCol)
	pConditionsCount := pbeam.Count(scope, pConditions, pbeam.CountParams{
		PartitionSelectionParams: pbeam.PartitionSelectionParams{
			Epsilon: budgetShare.PartitionEpsilon,
			Delta:   budgetShare.PartitionDelta,
		},
		AggregationEpsilon:       budgetShare.AggregationEpsilon,
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

func CountTestResultsDp(scope beam.Scope, col beam.PCollection, pSpec *pbeam.PrivacySpec, budgetShare healtcaredp.DpBudgetShare) beam.PCollection {
	scope = scope.Scope("CountTestResultsDp")
	pCol := pbeam.MakePrivateFromStruct(scope, col, pSpec, "Name")

	pTestResults := pbeam.ParDo(scope, extractTestResult, pCol)
	pTestResultsCount := pbeam.Count(scope, pTestResults, pbeam.CountParams{
		PartitionSelectionParams: pbeam.PartitionSelectionParams{
			Epsilon: budgetShare.PartitionEpsilon,
			Delta:   budgetShare.PartitionDelta,
		},
		AggregationEpsilon:       budgetShare.AggregationEpsilon,
		MaxPartitionsContributed: 8,
		MaxValue:                 24,
	})
	return pTestResultsCount
}
