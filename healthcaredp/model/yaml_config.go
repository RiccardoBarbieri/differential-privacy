package model

import (
	"fmt"
	"github.com/goccy/go-yaml"
	log "github.com/golang/glog"
	"os"
	"strings"
)

type OperationType struct {
	OperationName string  `yaml:"name"`
	OperationType string  `yaml:"type"`
	Column        string  `yaml:"column"`
	GenerateNonDp bool    `yaml:"generate_non_dp"`
	Importance    float64 `yaml:"importance"`
}

type ConfigurationType struct {
	DataDir        string `yaml:"data_dir"`
	Input          string `yaml:"input"`
	OutputBaseName string `yaml:"output_base_name"`
	IdField        string `yaml:"id_field"`
}

type TypeType struct {
	Column string `yaml:"column"`
	Type   string `yaml:"type"`
}

type PrivacyBudgetType struct {
	Delta            float64 `yaml:"delta"`
	Epsilon          float64 `yaml:"epsilon"`
	AggregationShare float64 `yaml:"aggregation_share"`
}

type PipelineDp struct {
	Configuration ConfigurationType `yaml:"configuration"`
	PrivacyBudget PrivacyBudgetType `yaml:"privacy_budget"`
	Types         []TypeType        `yaml:"types"`
	Operations    []OperationType   `yaml:"operations"`
}

type YamlConfig struct {
	PipelineDp PipelineDp `yaml:"pipelinedp"`
}

func LoadYamlConfig(filename string) (config *YamlConfig, err error) {
	log.Infof("Loading config from file: %s", filename)
	var yaml_file []byte
	yaml_file, err = os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	config = new(YamlConfig)
	err = yaml.Unmarshal(yaml_file, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	for i, typ := range config.PipelineDp.Types {
		cleanHeader := strings.ReplaceAll(typ.Column, " ", "")
		config.PipelineDp.Types[i].Column = cleanHeader
	}
	for i, op := range config.PipelineDp.Operations {
		cleanColumn := strings.ReplaceAll(op.Column, " ", "")
		config.PipelineDp.Operations[i].Column = cleanColumn
	}
	return config, nil
}
