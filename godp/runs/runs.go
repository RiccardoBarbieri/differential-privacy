package runs

import (
	"context"
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/runners/direct"
	"github.com/spf13/cobra"
	"godp"
	"godp/aggregations"
	"godp/model"
	utils2 "godp/model/utils"
	"godp/utils"
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

	RunConditionsCount(healthcaredp.GlobalScope,
		healthcaredp.CurrentIOArgs.OutputCsv,
		healthcaredp.CurrentIOArgs.GenerateNonDp,
		healthcaredp.AdmissionsCleaned)
	RunTestResultsCount(healthcaredp.GlobalScope,
		healthcaredp.CurrentIOArgs.OutputCsv,
		healthcaredp.CurrentIOArgs.GenerateNonDp,
		healthcaredp.AdmissionsCleaned)

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), healthcaredp.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	headers, err := utils2.StructCsvHeaders(model.Admission{})
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	utils2.WriteHeaders(healthcaredp.CurrentIOArgs.OutputClean, headers...)
	ConditionsCountWriteHeaders(healthcaredp.CurrentIOArgs.GenerateNonDp)
	TestResultsCountWriteHeaders(healthcaredp.CurrentIOArgs.GenerateNonDp)

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
		RunConditionsCount(healthcaredp.GlobalScope,
			healthcaredp.CurrentIOArgs.OutputCsv,
			healthcaredp.CurrentIOArgs.GenerateNonDp,
			healthcaredp.AdmissionsCleaned)
	case "CountTestResults":
		RunTestResultsCount(healthcaredp.GlobalScope,
			healthcaredp.CurrentIOArgs.OutputCsv,
			healthcaredp.CurrentIOArgs.GenerateNonDp,
			healthcaredp.AdmissionsCleaned)
	}

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), healthcaredp.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	headers, err := utils2.StructCsvHeaders(model.Admission{})
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	utils2.WriteHeaders(healthcaredp.CurrentIOArgs.OutputClean, headers...)
	switch args[0] {
	case "CountConditions":
		ConditionsCountWriteHeaders(healthcaredp.CurrentIOArgs.GenerateNonDp)
	case "CountTestResults":
		TestResultsCountWriteHeaders(healthcaredp.CurrentIOArgs.GenerateNonDp)
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
		MeanStayByWeek(healthcaredp.GlobalScope,
			healthcaredp.CurrentIOArgs.OutputCsv,
			healthcaredp.CurrentIOArgs.GenerateNonDp,
			healthcaredp.AdmissionsCleaned)
	}

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), healthcaredp.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	headers, err := utils2.StructCsvHeaders(model.Admission{})
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	utils2.WriteHeaders(healthcaredp.CurrentIOArgs.OutputClean, headers...)
	switch args[0] {
	case "MeanStayByWeek":
		MeanStayByWeekWriteHeaders(healthcaredp.CurrentIOArgs.GenerateNonDp)
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
	healthcaredp.GlobalPipeline = beam.NewPipeline()
	healthcaredp.GlobalScope = healthcaredp.GlobalPipeline.Root()

	var datasetFilename = config.PipelineDp.Configuration.DataDir + "/" + config.PipelineDp.Configuration.Input

	model.Headers, err = utils2.GetHeaders(datasetFilename)
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	for i, val := range model.Headers {
		if val == config.PipelineDp.Configuration.IdField {
			model.IdFieldIndex = i
			break
		}
	}
	newDFilename, err := utils2.RemoveHeadersAndSaveCsv(datasetFilename)
	if err != nil {
		return fmt.Errorf("error removing headers: %v", err)
	}
	model.TypesMap, err = model.CompileTypesMap(config.PipelineDp.Types)
	if err != nil {
		return fmt.Errorf("error compiling types map: %v", err)
	}
	pcol := utils2.ReadGenericInput(healthcaredp.GlobalScope, newDFilename)

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
			utils2.WriteOutput(healthcaredp.GlobalScope, *pColCount, outFilename)
			utils2.PrintConsole(healthcaredp.GlobalScope, *pColCount)
		case "mean_per_key":
			pColMean, err := aggregations.MeanColumnByKey(healthcaredp.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating mean: %v", err)
			}
			utils2.WriteOutput(healthcaredp.GlobalScope, *pColMean, outFilename)
			utils2.PrintConsole(healthcaredp.GlobalScope, *pColMean)
		case "sum_per_key":
			pColSum, err := aggregations.SumColumnByKey(healthcaredp.GlobalScope, pcol, op, healthcaredp.Budget)
			if err != nil {
				return fmt.Errorf("error calculating sum: %v", err)
			}
			utils2.WriteOutput(healthcaredp.GlobalScope, *pColSum, outFilename)
			utils2.PrintConsole(healthcaredp.GlobalScope, *pColSum)
		}
	}

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), healthcaredp.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	return nil
}
