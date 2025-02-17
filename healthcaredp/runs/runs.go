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

	pColCount := aggregations.CountColumn(model.GlobalScope, pcol, "TestResults", healthcaredp.Budget)

	model.PrintConsole(model.GlobalScope, pColCount)

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), model.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	return nil
}
