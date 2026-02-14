package healthcaredp

import (
	"godp/model"
	"math"

	"github.com/google/differential-privacy/privacy-on-beam/v3/pbeam"
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

func (db DpBudget) GetBudgetShare(operation string) DpBudgetShare {
	return db.BudgetShares[operation]
}
