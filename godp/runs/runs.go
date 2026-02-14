package runs

import (
	"context"
	"fmt"
	"godp"
	"godp/aggregations"
	"godp/model"
	modelutils "godp/model/utils"
	"godp/utils"
	"strings"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/runners/direct"
	"github.com/spf13/cobra"
)

func RunFromFile(cmd *cobra.Command, args []string) (err error) {
	var config *model.YamlConfig
	var filename string

	filename, err = cmd.PersistentFlags().GetString("file")
	if err != nil {
		return fmt.Errorf("error getting config file parameter: %v", err)
	}

	generateNonDpFlag, err := cmd.Flags().GetBool("generate-non-dp")
	if err != nil {
		return fmt.Errorf("error getting generate-non-dp parameter: %v", err)
	}

	printConsoleFlag, err := cmd.Flags().GetBool("print-console")
	if err != nil {
		return fmt.Errorf("error getting print-console parameter: %v", err)
	}

	config, err = model.LoadYamlConfig(filename)
	if err != nil {
		return fmt.Errorf("error loading config file: %v", err)
	}

	err = healthcaredp.Budget.InitYamlBudgetShares(config)
	if err != nil {
		return err
	}

	beam.Init()
	healthcaredp.GlobalPipeline = beam.NewPipeline()
	healthcaredp.GlobalScope = healthcaredp.GlobalPipeline.Root()

	var datasetFilename = config.PipelineDp.Configuration.DataDir + "/" + config.PipelineDp.Configuration.Input

	model.Headers, err = modelutils.GetHeaders(datasetFilename)
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	for i, val := range model.Headers {
		if val == config.PipelineDp.Configuration.IdField {
			model.IdFieldIndex = i
			break
		}
	}
	newDFilename, err := modelutils.RemoveHeadersAndSaveCsv(datasetFilename)
	if err != nil {
		return fmt.Errorf("error removing headers: %v", err)
	}
	model.TypesMap, err = model.CompileTypesMap(config.PipelineDp.Types)
	if err != nil {
		return fmt.Errorf("error compiling types map: %v", err)
	}
	pcol := modelutils.ReadGenericInput(healthcaredp.GlobalScope, newDFilename)

	for _, op := range config.PipelineDp.Operations {
		if !utils.SliceContains(model.Headers, op.Column) {
			return fmt.Errorf("column %s does not exists in the dataset", op.Column)
		}

		outFilename := fmt.Sprintf("%s/%s_%s.csv",
			config.PipelineDp.Configuration.DataDir,
			strings.TrimSuffix(config.PipelineDp.Configuration.OutputBaseName, ".csv"),
			op.OperationName)

		if generateNonDpFlag {
			outFilenameClear := fmt.Sprintf("%s/%s_%s_clear.csv",
				config.PipelineDp.Configuration.DataDir,
				strings.TrimSuffix(config.PipelineDp.Configuration.OutputBaseName, ".csv"),
				op.OperationName)

			switch op.OperationType {
			case "count":
				pColCountClear, err := aggregations.CountColumnClear(healthcaredp.GlobalScope, pcol, op)
				if err != nil {
					return fmt.Errorf("error calculating clear count: %v", err)
				}
				modelutils.WriteOutput(healthcaredp.GlobalScope, *pColCountClear, outFilenameClear)
				if printConsoleFlag {
					modelutils.PrintConsole(healthcaredp.GlobalScope, *pColCountClear)
				}
			case "mean_per_key":
				pColMeanClear, err := aggregations.MeanColumnByKeyClear(healthcaredp.GlobalScope, pcol, op)
				if err != nil {
					return fmt.Errorf("error calculating clear mean: %v", err)
				}
				modelutils.WriteOutput(healthcaredp.GlobalScope, *pColMeanClear, outFilenameClear)
				if printConsoleFlag {
					modelutils.PrintConsole(healthcaredp.GlobalScope, *pColMeanClear)
				}
			case "sum_per_key":
				pColSumClear, err := aggregations.SumColumnByKeyClear(healthcaredp.GlobalScope, pcol, op)
				if err != nil {
					return fmt.Errorf("error calculating clear sum: %v", err)
				}
				modelutils.WriteOutput(healthcaredp.GlobalScope, *pColSumClear, outFilenameClear)
				if printConsoleFlag {
					modelutils.PrintConsole(healthcaredp.GlobalScope, *pColSumClear)
				}
			}
			fmt.Printf("Results for %s operation (clear) saved in file %s\n\n", op.OperationName, outFilenameClear)
		}

		var pColRes *beam.PCollection
		switch op.OperationType {
		case "count":
			pColRes, err = aggregations.CountColumn(healthcaredp.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating count: %v", err)
			}
		case "mean_per_key":
			pColRes, err = aggregations.MeanColumnByKey(healthcaredp.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating mean: %v", err)
			}
		case "sum_per_key":
			pColRes, err = aggregations.SumColumnByKey(healthcaredp.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating sum: %v", err)
			}
		default:
			return fmt.Errorf("operation type %s not supported", op.OperationType)
		}
		modelutils.WriteOutput(healthcaredp.GlobalScope, *pColRes, outFilename)
		if printConsoleFlag {
			modelutils.PrintConsole(healthcaredp.GlobalScope, *pColRes)
		}

		fmt.Printf("Results for %s operation saved in file %s\n\n", op.OperationName, outFilename)
	}

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), healthcaredp.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	//_ = modelutils.DeleteFile(newDFilename)

	return nil
}
