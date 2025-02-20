package runs

import (
	"context"
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/runners/direct"
	"github.com/spf13/cobra"
	"healthcaredp"
	"healthcaredp/aggregations"
	"healthcaredp/model"
	"healthcaredp/utils"
	"strings"
)

func RunAll(cmd *cobra.Command, args []string) (err error) {
	err = healthcaredp.Budget.InitAllBudgetShares(
		map[string]float64{
			"CountConditions":  1.0,
			"CountTestResults": 1.0,
			"MeanStayByWeek":   1.0,
		},
	)
	if err != nil {
		return err
	}

	RunConditionsCount(model.GlobalScope,
		model.CurrentIOArgs.OutputCsv,
		model.CurrentIOArgs.GenerateNonDp,
		model.AdmissionsCleaned)
	RunTestResultsCount(model.GlobalScope,
		model.CurrentIOArgs.OutputCsv,
		model.CurrentIOArgs.GenerateNonDp,
		model.AdmissionsCleaned)

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), model.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	headers, err := model.StructCsvHeaders(model.Admission{})
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	model.WriteHeaders(model.CurrentIOArgs.OutputClean, headers...)
	ConditionsCountWriteHeaders(model.CurrentIOArgs.GenerateNonDp)
	TestResultsCountWriteHeaders(model.CurrentIOArgs.GenerateNonDp)

	return nil
}

func RunCounts(cmd *cobra.Command, args []string) (err error) {
	err = healthcaredp.Budget.InitBudgetShares(
		map[string]float64{
			"CountConditions":  1.0,
			"CountTestResults": 1.0,
		},
	)
	if err != nil {
		return err
	}

	switch args[0] {
	case "CountConditions":
		RunConditionsCount(model.GlobalScope,
			model.CurrentIOArgs.OutputCsv,
			model.CurrentIOArgs.GenerateNonDp,
			model.AdmissionsCleaned)
	case "CountTestResults":
		RunTestResultsCount(model.GlobalScope,
			model.CurrentIOArgs.OutputCsv,
			model.CurrentIOArgs.GenerateNonDp,
			model.AdmissionsCleaned)
	}

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), model.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	headers, err := model.StructCsvHeaders(model.Admission{})
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	model.WriteHeaders(model.CurrentIOArgs.OutputClean, headers...)
	switch args[0] {
	case "CountConditions":
		ConditionsCountWriteHeaders(model.CurrentIOArgs.GenerateNonDp)
	case "CountTestResults":
		TestResultsCountWriteHeaders(model.CurrentIOArgs.GenerateNonDp)
	}

	return nil
}

func RunMeans(cmd *cobra.Command, args []string) (err error) {
	err = healthcaredp.Budget.InitBudgetShares(
		map[string]float64{
			"MeanStayByWeek": 1.0,
		},
	)
	if err != nil {
		return err
	}

	switch args[0] {
	case "MeanStayByWeek":
		MeanStayByWeek(model.GlobalScope,
			model.CurrentIOArgs.OutputCsv,
			model.CurrentIOArgs.GenerateNonDp,
			model.AdmissionsCleaned)
	}

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), model.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	headers, err := model.StructCsvHeaders(model.Admission{})
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	model.WriteHeaders(model.CurrentIOArgs.OutputClean, headers...)
	switch args[0] {
	case "MeanStayByWeek":
		MeanStayByWeekWriteHeaders(model.CurrentIOArgs.GenerateNonDp)
	}

	return nil
}

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
	model.GlobalPipeline = beam.NewPipeline()
	model.GlobalScope = model.GlobalPipeline.Root()

	var datasetFilename = config.PipelineDp.Configuration.DataDir + "/" + config.PipelineDp.Configuration.Input

	model.Headers, err = model.GetHeaders(datasetFilename)
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	for i, val := range model.Headers {
		if val == config.PipelineDp.Configuration.IdField {
			model.IdFieldIndex = i
			break
		}
	}
	newDFilename, err := model.RemoveHeadersAndSaveCsv(datasetFilename)
	if err != nil {
		return fmt.Errorf("error removing headers: %v", err)
	}
	model.TypesMap, err = model.CompileTypesMap(config.PipelineDp.Types)
	if err != nil {
		return fmt.Errorf("error compiling types map: %v", err)
	}
	pcol := model.ReadGenericInput(model.GlobalScope, newDFilename)

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
			pColCount, err := aggregations.CountColumn(model.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating count: %v", err)
			}
			model.WriteOutput(model.GlobalScope, *pColCount, outFilename)
			model.PrintConsole(model.GlobalScope, *pColCount)
		case "mean_per_key":
			pColMean, err := aggregations.MeanColumnByKey(model.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating mean: %v", err)
			}
			model.WriteOutput(model.GlobalScope, *pColMean, outFilename)
			model.PrintConsole(model.GlobalScope, *pColMean)
		case "sum_per_key":
			pColSum, err := aggregations.SumColumnByKey(model.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating sum: %v", err)
			}
			model.WriteOutput(model.GlobalScope, *pColSum, outFilename)
			model.PrintConsole(model.GlobalScope, *pColSum)
		}
	}

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), model.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	return nil
}
