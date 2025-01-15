package healthcaredp

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/google/differential-privacy/privacy-on-beam/v3/pbeam"
	"healthcaredp/utils"
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

var SupportedOperations []string

var Budget DpBudget

const delta = 1e-5

var epsilon = math.Log(3)
var aggregationEpsilon = epsilon / 2
var partitionEpsilon = epsilon - aggregationEpsilon

var CountOperations = []string{"CountConditions", "CountTestResults"}
var MeanOperations = []string{"MeanStayByWeek"}

func init() {

	pSpecParams := pbeam.PrivacySpecParams{
		AggregationEpsilon:        aggregationEpsilon,
		PartitionSelectionEpsilon: partitionEpsilon,
		PartitionSelectionDelta:   delta,
	}

	SupportedOperations = append(SupportedOperations, CountOperations...)
	SupportedOperations = append(SupportedOperations, MeanOperations...)

	var err error
	Budget.Delta = delta
	Budget.Epsilon = epsilon
	Budget.PrivacySpec, err = pbeam.NewPrivacySpec(pSpecParams)
	Budget.BudgetShares = make(map[string]DpBudgetShare)
	if err != nil {
		log.Fatalf("Failed to create privacy spec: %v", err)
	}
}

func (db DpBudget) InitAllBudgetShares(importance map[string]float64) (err error) {
	if len(importance) != len(SupportedOperations) {
		return fmt.Errorf("expected %d importance values, got %d", len(SupportedOperations), len(importance))
	}
	for _, key := range SupportedOperations {
		if _, ok := importance[key]; !ok {
			return fmt.Errorf("importance map missing operation: %s", key)
		}
	}
	totalImportance := 0.0
	for _, imp := range importance {
		totalImportance += imp
	}
	for key, imp := range importance {
		db.BudgetShares[key] = DpBudgetShare{
			AggregationEpsilon: (imp / totalImportance) * aggregationEpsilon,
			PartitionEpsilon:   (imp / totalImportance) * partitionEpsilon,
			AggregationDelta:   (imp / totalImportance) * delta,
			PartitionDelta:     delta - (imp/totalImportance)*delta,
		}
	}
	return nil
}

func (db DpBudget) InitBudgetShares(importance map[string]float64) (err error) {
	if len(importance) > len(SupportedOperations) {
		return fmt.Errorf("importance map has more than %d operations", len(SupportedOperations))
	}
	for key, _ := range importance {
		if !utils.SliceContains(SupportedOperations, key) {
			return fmt.Errorf("importance map contains unsupported operation: %s", key)
		}
	}
	totalImportance := 0.0
	for _, imp := range importance {
		totalImportance += imp
	}
	for key, imp := range importance {
		db.BudgetShares[key] = DpBudgetShare{
			AggregationEpsilon: (imp / totalImportance) * aggregationEpsilon,
			PartitionEpsilon:   (imp / totalImportance) * partitionEpsilon,
			AggregationDelta:   (imp / totalImportance) * delta,
			PartitionDelta:     delta - (imp/totalImportance)*delta,
		}
	}
	return nil
}

func (db DpBudget) GetBudgetShare(operation string) DpBudgetShare {
	return db.BudgetShares[operation]
}
