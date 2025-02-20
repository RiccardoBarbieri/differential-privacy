package healthcaredp

import (
	"fmt"
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

const GAUSSIAN = "gauss"
const LAPLACIAN = "laplace"

var SupportedOperations []string

var Budget DpBudget

var Delta = 1e-5

var Epsilon = math.Log(3)
var AggregationEpsilon = Epsilon / 2
var PartitionEpsilon = Epsilon - AggregationEpsilon

var CountOperations = []string{"CountConditions", "CountTestResults"}
var MeanOperations = []string{"MeanStayByWeek"}

func init() {

	SupportedOperations = append(SupportedOperations, CountOperations...)
	SupportedOperations = append(SupportedOperations, MeanOperations...)
}

func (db *DpBudget) InitYamlBudgetShares(config *model.YamlConfig) (err error) {
	var delta = config.PipelineDp.PrivacyBudget.Delta
	var epsilon = config.PipelineDp.PrivacyBudget.Epsilon

	var aggregationEpsilon = epsilon * config.PipelineDp.PrivacyBudget.AggregationShare

	var aggregationDelta = 0.0
	var partitionSelectionDelta = 0.0
	if config.PipelineDp.PrivacyBudget.NoiseKind == GAUSSIAN {
		aggregationDelta = delta * config.PipelineDp.PrivacyBudget.AggregationShare
		partitionSelectionDelta = delta - aggregationDelta
	} else if config.PipelineDp.PrivacyBudget.NoiseKind == LAPLACIAN {
		partitionSelectionDelta = delta
	}
	var partitionSelectionEpsilon = epsilon - aggregationEpsilon

	var pSpecParams = pbeam.PrivacySpecParams{}
	if config.PipelineDp.PrivacyBudget.NoiseKind == GAUSSIAN {
		pSpecParams = pbeam.PrivacySpecParams{
			AggregationEpsilon:        aggregationEpsilon,
			AggregationDelta:          aggregationDelta,
			PartitionSelectionEpsilon: partitionSelectionEpsilon,
			PartitionSelectionDelta:   partitionSelectionDelta,
		}
	} else if config.PipelineDp.PrivacyBudget.NoiseKind == LAPLACIAN {
		pSpecParams = pbeam.PrivacySpecParams{
			AggregationEpsilon:        aggregationEpsilon,
			PartitionSelectionEpsilon: partitionSelectionEpsilon,
			PartitionSelectionDelta:   delta,
		}
	}

	db.Delta = delta
	db.Epsilon = epsilon
	db.PrivacySpec, err = pbeam.NewPrivacySpec(pSpecParams)
	if err != nil {
		return err
	}
	db.BudgetShares = make(map[string]DpBudgetShare)

	totalImportance := 0.0
	for _, op := range config.PipelineDp.Operations {
		totalImportance += op.Importance
	}
	for _, op := range config.PipelineDp.Operations {
		if config.PipelineDp.PrivacyBudget.NoiseKind == GAUSSIAN {
			db.BudgetShares[op.OperationName] = DpBudgetShare{
				AggregationEpsilon: (op.Importance / totalImportance) * aggregationEpsilon,
				PartitionEpsilon:   (op.Importance / totalImportance) * partitionSelectionEpsilon,
				AggregationDelta:   (op.Importance / totalImportance) * aggregationDelta,
				PartitionDelta:     (op.Importance / totalImportance) * partitionSelectionDelta,
			}
		} else if config.PipelineDp.PrivacyBudget.NoiseKind == LAPLACIAN {
			db.BudgetShares[op.OperationName] = DpBudgetShare{
				AggregationEpsilon: (op.Importance / totalImportance) * aggregationEpsilon,
				PartitionEpsilon:   (op.Importance / totalImportance) * partitionSelectionEpsilon,
				PartitionDelta:     (op.Importance / totalImportance) * delta,
			}
		}
	}
	return nil
}

func (db DpBudget) InitAllBudgetShares(importance map[string]float64) (err error) {
	pSpecParams := pbeam.PrivacySpecParams{
		AggregationEpsilon:        AggregationEpsilon,
		PartitionSelectionEpsilon: PartitionEpsilon,
		PartitionSelectionDelta:   Delta,
	}
	db.Delta = Delta
	db.Epsilon = Epsilon
	db.PrivacySpec, err = pbeam.NewPrivacySpec(pSpecParams)
	db.BudgetShares = make(map[string]DpBudgetShare)
	if err != nil {
		return fmt.Errorf("failed to create privacy spec: %v", err)
	}

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
	pSpecParams := pbeam.PrivacySpecParams{
		AggregationEpsilon:        AggregationEpsilon,
		PartitionSelectionEpsilon: PartitionEpsilon,
		PartitionSelectionDelta:   Delta,
	}
	db.Delta = Delta
	db.Epsilon = Epsilon
	db.PrivacySpec, err = pbeam.NewPrivacySpec(pSpecParams)
	db.BudgetShares = make(map[string]DpBudgetShare)
	if err != nil {
		return fmt.Errorf("failed to create privacy spec: %v", err)
	}

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
