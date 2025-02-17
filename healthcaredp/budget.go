package healthcaredp

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/google/differential-privacy/privacy-on-beam/v3/pbeam"
	"healthcaredp/model"
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

var Delta = 1e-5

var Epsilon = math.Log(3)
var AggregationEpsilon = Epsilon / 2
var PartitionEpsilon = Epsilon - AggregationEpsilon

var CountOperations = []string{"CountConditions", "CountTestResults"}
var MeanOperations = []string{"MeanStayByWeek"}

func init() {

	pSpecParams := pbeam.PrivacySpecParams{
		AggregationEpsilon:        AggregationEpsilon,
		PartitionSelectionEpsilon: PartitionEpsilon,
		PartitionSelectionDelta:   Delta,
	}

	SupportedOperations = append(SupportedOperations, CountOperations...)
	SupportedOperations = append(SupportedOperations, MeanOperations...)

	var err error
	Budget.Delta = Delta
	Budget.Epsilon = Epsilon
	Budget.PrivacySpec, err = pbeam.NewPrivacySpec(pSpecParams)
	Budget.BudgetShares = make(map[string]DpBudgetShare)
	if err != nil {
		log.Fatalf("Failed to create privacy spec: %v", err)
	}
}

func (db DpBudget) InitYamlBudgetShares(config *model.YamlConfig) (err error) {

	var delta = config.PipelineDp.PrivacyBudget.Delta
	var epsilon = config.PipelineDp.PrivacyBudget.Epsilon

	pSpecParams := pbeam.PrivacySpecParams{
		AggregationEpsilon:        epsilon * config.PipelineDp.PrivacyBudget.AggregationShare,
		PartitionSelectionEpsilon: epsilon - epsilon*config.PipelineDp.PrivacyBudget.AggregationShare,
		PartitionSelectionDelta:   delta,
	}

	Budget.Delta = config.PipelineDp.PrivacyBudget.Delta
	Budget.Epsilon = config.PipelineDp.PrivacyBudget.Epsilon
	Budget.PrivacySpec, err = pbeam.NewPrivacySpec(pSpecParams)
	if err != nil {
		return err
	}
	Budget.BudgetShares = make(map[string]DpBudgetShare)

	totalImportance := 0.0
	for _, op := range config.PipelineDp.Operations {
		totalImportance += op.Importance
	}
	for _, op := range config.PipelineDp.Operations {
		db.BudgetShares[op.OperationName] = DpBudgetShare{
			AggregationEpsilon: (op.Importance / totalImportance) * AggregationEpsilon,
			PartitionEpsilon:   (op.Importance / totalImportance) * PartitionEpsilon,
			AggregationDelta:   (op.Importance / totalImportance) * Delta,
			PartitionDelta:     Delta - (op.Importance/totalImportance)*Delta,
		}
	}
	return nil
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
			AggregationEpsilon: (imp / totalImportance) * AggregationEpsilon,
			PartitionEpsilon:   (imp / totalImportance) * PartitionEpsilon,
			AggregationDelta:   (imp / totalImportance) * Delta,
			PartitionDelta:     Delta - (imp/totalImportance)*Delta,
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
			AggregationEpsilon: (imp / totalImportance) * AggregationEpsilon,
			PartitionEpsilon:   (imp / totalImportance) * PartitionEpsilon,
			AggregationDelta:   (imp / totalImportance) * Delta,
			PartitionDelta:     Delta - (imp/totalImportance)*Delta,
		}
	}
	return nil
}

func (db DpBudget) GetBudgetShare(operation string) DpBudgetShare {
	return db.BudgetShares[operation]
}
