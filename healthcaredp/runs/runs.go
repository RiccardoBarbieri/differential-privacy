package runs

import (
	"context"
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/runners/direct"
	"github.com/spf13/cobra"
	"healthcaredp"
	"healthcaredp/model"
	"healthcaredp/utils"
)

func RunAll(cmd *cobra.Command, args []string) (err error) {

	err = healthcaredp.Budget.InitAllBudgetShares(
		map[string]float64{
			"CountConditions":  1.0,
			"CountTestResults": 1.0,
			"AvgStayByWeek":    1.0,
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

	headers, err := utils.StructCsvHeaders(model.Admission{})
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	utils.WriteHeaders(healthcaredp.CurrentIOArgs.OutputClean, headers...)
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

	headers, err := utils.StructCsvHeaders(model.Admission{})
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	utils.WriteHeaders(healthcaredp.CurrentIOArgs.OutputClean, headers...)
	switch args[0] {
	case "CountConditions":
		ConditionsCountWriteHeaders(healthcaredp.CurrentIOArgs.GenerateNonDp)
	case "CountTestResults":
		TestResultsCountWriteHeaders(healthcaredp.CurrentIOArgs.GenerateNonDp)
	}

	return nil
}

func RunAvg(cmd *cobra.Command, args []string) (err error) {

	err = healthcaredp.Budget.InitBudgetShares(
		map[string]float64{
			"AvgStayByWeek": 1.0,
		},
	)
	if err != nil {
		return err
	}

	switch args[0] {
	case "AvgStayByWeek":
		RunAvgStayByWeek(healthcaredp.GlobalScope,
			healthcaredp.CurrentIOArgs.OutputCsv,
			healthcaredp.CurrentIOArgs.GenerateNonDp,
			healthcaredp.AdmissionsCleaned)
	}

	// Execute pipeline.
	_, err = direct.Execute(context.Background(), healthcaredp.GlobalPipeline)
	if err != nil {
		return fmt.Errorf("error executing pipeline: %v", err)
	}

	headers, err := utils.StructCsvHeaders(model.Admission{})
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}
	utils.WriteHeaders(healthcaredp.CurrentIOArgs.OutputClean, headers...)
	switch args[0] {
	case "AvgStayByWeek":
		AvgStayByWeekWriteHeaders(healthcaredp.CurrentIOArgs.GenerateNonDp)
	}

	return nil
}
