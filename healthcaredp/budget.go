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

type DpBudget struct {
	PrivacySpec  *pbeam.PrivacySpec
	BudgetShares map[string]DpBudgetShare
	Delta        float64
	Epsilon      float64
}

const delta = 1e-5

var epsilon = math.Log(3)
var aggregationEpsilon = epsilon / 2
var partitionEpsilon = epsilon - aggregationEpsilon

var supportedOperations = []string{"CountConditionsDp", "CountTestResultsDp"}

var Budget DpBudget

func init() {
	pSpecParams := pbeam.PrivacySpecParams{
		AggregationEpsilon:        aggregationEpsilon,
		PartitionSelectionEpsilon: partitionEpsilon,
		PartitionSelectionDelta:   delta,
	}
	var err error
	Budget.Delta = delta
	Budget.Epsilon = epsilon
	Budget.PrivacySpec, err = pbeam.NewPrivacySpec(pSpecParams)
	Budget.BudgetShares = make(map[string]DpBudgetShare)
	if err != nil {
		log.Fatalf("Failed to create privacy spec: %v", err)
	}
}

func (db DpBudget) InitAllBudgetShares(importance ...float64) {
	if len(importance) != len(supportedOperations) {
		log.Fatalf("Expected %d importance values, got %d", len(supportedOperations), len(importance))
	}
	totalImportance := 0.0
	for _, imp := range importance {
		totalImportance += imp
	}
	for i, imp := range importance {
		db.BudgetShares[supportedOperations[i]] = DpBudgetShare{
			AggregationEpsilon: (imp / totalImportance) * aggregationEpsilon,
			PartitionEpsilon:   (imp / totalImportance) * partitionEpsilon,
			AggregationDelta:   (imp / totalImportance) * delta,
			PartitionDelta:     delta - (imp/totalImportance)*delta,
		}
	}
}

func (db DpBudget) GetBudgetShare(operation string) DpBudgetShare {
	return db.BudgetShares[operation]
}
