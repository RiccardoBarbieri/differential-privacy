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

		fmt.Printf("Results for %s operation saved in file %s\n\n", op.OperationName, outFilename)
		switch op.OperationType {
		case "count":
			pColCount, err := aggregations.CountColumn(healthcaredp.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating count: %v", err)
			}
			modelutils.WriteOutput(healthcaredp.GlobalScope, *pColCount, outFilename)
			modelutils.PrintConsole(healthcaredp.GlobalScope, *pColCount)
		case "mean_per_key":
			pColMean, err := aggregations.MeanColumnByKey(healthcaredp.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating mean: %v", err)
			}
			modelutils.WriteOutput(healthcaredp.GlobalScope, *pColMean, outFilename)
			modelutils.PrintConsole(healthcaredp.GlobalScope, *pColMean)
		case "sum_per_key":
			pColSum, err := aggregations.SumColumnByKey(healthcaredp.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating sum: %v", err)
			}
			modelutils.WriteOutput(healthcaredp.GlobalScope, *pColSum, outFilename)
			modelutils.PrintConsole(healthcaredp.GlobalScope, *pColSum)
		}
	}

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), healthcaredp.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	return nil
}
