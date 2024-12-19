package healtcaredp

import (
	log "github.com/golang/glog"
	"github.com/google/differential-privacy/privacy-on-beam/v3/pbeam"
	"math"
)

type DpBudgetShare struct {
	AggregationEpsilon float64
	PartitionEpsilon   float64
	AggregationDelta   float64
	PartitionDelta     float64
}

type budgetSplit struct {
	aggregationEpsilon float64
	partitionEpsilon   float64
	delta              float64
}

const delta = 1e-5

var epsilon = math.Log(3)
var aggregationEpsilon = epsilon / 2
var partitionEpsilon = epsilon - aggregationEpsilon

var GlobalPrivacySpec *pbeam.PrivacySpec

var supportedOperations = []string{"CountConditionsDp", "CountTestResultsDp"}
var budgetSplits = map[string]budgetSplit{}

func init() {
	pSpecParams := pbeam.PrivacySpecParams{
		AggregationEpsilon:        aggregationEpsilon,
		PartitionSelectionEpsilon: partitionEpsilon,
		PartitionSelectionDelta:   delta,
	}
	var err error
	GlobalPrivacySpec, err = pbeam.NewPrivacySpec(pSpecParams)
	if err != nil {
		log.Fatalf("Failed to create privacy spec: %v", err)
	}

}

func InitBudgetSplits(importance ...float64) {
	if len(importance) != len(supportedOperations) {
		log.Fatalf("Expected %d importance values, got %d", len(supportedOperations), len(importance))
	}
	totalImportance := 0.0
	for _, imp := range importance {
		totalImportance += imp
	}
	for i, imp := range importance {
		budgetSplits[supportedOperations[i]] = budgetSplit{
			aggregationEpsilon: (imp / totalImportance) * aggregationEpsilon,
			partitionEpsilon:   (imp / totalImportance) * partitionEpsilon,
			delta:              (imp / totalImportance) * delta,
		}
	}
}

func GetBudgetShare(operation string) DpBudgetShare {
	switch operation {
	case "CountConditionsDp":
		return DpBudgetShare{
			AggregationEpsilon: budgetSplits["CountConditionsDp"].aggregationEpsilon,
			PartitionEpsilon:   budgetSplits["CountConditionsDp"].partitionEpsilon,

			AggregationDelta: budgetSplits["CountConditionsDp"].delta,
			PartitionDelta:   delta - budgetSplits["CountConditionsDp"].delta,
		}
	case "CountTestResultsDp":
		return DpBudgetShare{
			AggregationEpsilon: budgetSplits["CountTestResultsDp"].aggregationEpsilon,
			PartitionEpsilon:   budgetSplits["CountTestResultsDp"].partitionEpsilon,

			AggregationDelta: budgetSplits["CountTestResultsDp"].delta,
			PartitionDelta:   delta - budgetSplits["CountTestResultsDp"].delta,
		}
	default:
		log.Fatalf("Unsupported operation: %s", operation)
		return DpBudgetShare{}
	}
}
