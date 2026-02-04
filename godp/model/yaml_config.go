package model

import (
	"fmt"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"

	"os"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	log "github.com/golang/glog"
)

type OperationType struct {
	OperationName string            `yaml:"name" validate:"required"`
	OperationType string            `yaml:"type" validate:"required"`
	Column        string            `yaml:"column" validate:"required"`
	KeyColumn     *string           `yaml:"key_column,omitempty"`
	Importance    float64           `yaml:"importance" validate:"required"`
	PrivacyParams PrivacyParamsType `yaml:"privacy_params" validate:"required"`
}

type PrivacyParamsType struct {
	MaxCategoriesContributed    *int64   `yaml:"max_categories_contributed,omitempty"`
	MaxContributions            *int64   `yaml:"max_contributions,omitempty"`
	MaxContributionsPerCategory *int64   `yaml:"max_contributions_per_category,omitempty"`
	MinValue                    *float64 `yaml:"min_value,omitempty"`
	MaxValue                    *float64 `yaml:"max_value,omitempty"`
}

type ConfigurationType struct {
	DataDir        string `yaml:"data_dir" validate:"required,dir"`
	Input          string `yaml:"input" validate:"required,endswith=csv"`
	OutputBaseName string `yaml:"output_base_name" validate:"required,endswith=csv"`
	IdField        string `yaml:"id_field" validate:"required"`
}

type TypeType struct {
	Column string `yaml:"column" validate:"required"`
	Type   string `yaml:"type" validate:"required,oneof=int string bool float date time datetime"`
}

type PrivacyBudgetType struct {
	NoiseKind        string  `yaml:"noise_kind" validate:"required,oneof=gauss laplace"`
	Delta            float64 `yaml:"delta" validate:"required"`
	Epsilon          float64 `yaml:"epsilon" validate:"required"`
	AggregationShare float64 `yaml:"aggregation_share" validate:"required,gt=0,lt=1"`
}

type PipelineDp struct {
	Configuration ConfigurationType `yaml:"configuration" validate:"required"`
	PrivacyBudget PrivacyBudgetType `yaml:"privacy_budget" validate:"required"`
	Types         []TypeType        `yaml:"types"`
	Operations    []OperationType   `yaml:"operations" validate:"required"`
}

type YamlConfig struct {
	PipelineDp PipelineDp `yaml:"pipelinedp" validate:"required"`
}

// Register custom error messages
func registerCustomErrorMessages(validate *validator.Validate, trans ut.Translator) (err error) {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("yaml")
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom error messages for specific tags
	err = validate.RegisterTranslation("required", trans,
		func(ut ut.Translator) error {
			return ut.Add("required", "{0} is a required field", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		},
	)
	if err != nil {
		return err
	}

	err = validate.RegisterTranslation("gt", trans, func(ut ut.Translator) error {
		return ut.Add("gt", "{0} must be greater than {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gt", fe.Field(), fe.Param())
		return t
	})
	if err != nil {
		return err
	}

	err = validate.RegisterTranslation("lt", trans, func(ut ut.Translator) error {
		return ut.Add("lt", "{0} must be less than {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("lt", fe.Field(), fe.Param())
		return t
	})
	if err != nil {
		return err
	}

	err = validate.RegisterTranslation("oneof", trans, func(ut ut.Translator) error {
		return ut.Add("oneof", "{0} must be one of {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("oneof", fe.Field(), fe.Param())
		return t
	})
	if err != nil {
		return err
	}

	err = validate.RegisterTranslation("endswith", trans, func(ut ut.Translator) error {
		return ut.Add("endswith", "{0} must end with {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("endswith", fe.Field(), fe.Param())
		return t
	})
	if err != nil {
		return err
	}

	return nil
}

func validateFieldsOpType(config *YamlConfig) error {
	for _, op := range config.PipelineDp.Operations {
		if op.OperationType == "mean_per_key" {
			if op.KeyColumn == nil {
				return fmt.Errorf("key_column is required for operation type: %s", op.OperationType)
			}
			if op.PrivacyParams.MaxCategoriesContributed == nil {
				return fmt.Errorf("max_categories_contributed is required for operation type: %s", op.OperationType)
			}
			if op.PrivacyParams.MaxContributionsPerCategory == nil {
				return fmt.Errorf("max_contributions_per_category is required for operation type: %s", op.OperationType)
			}
			if op.PrivacyParams.MinValue == nil {
				return fmt.Errorf("min_value is required for operation type: %s", op.OperationType)
			}
			if op.PrivacyParams.MaxValue == nil {
				return fmt.Errorf("max_value is required for operation type: %s", op.OperationType)
			}
		} else if op.OperationType == "count" {
			if op.PrivacyParams.MaxContributions == nil {
				return fmt.Errorf("max_value is required for operation type: %s", op.OperationType)
			}
			if op.PrivacyParams.MaxCategoriesContributed == nil {
				return fmt.Errorf("max_categories_contributed is required for operation type: %s", op.OperationType)
			}
		} else if op.OperationType == "sum_by_key" {
			if op.KeyColumn == nil {
				return fmt.Errorf("key_column is required for operation type: %s", op.OperationType)
			}
			if op.PrivacyParams.MaxCategoriesContributed == nil {
				return fmt.Errorf("max_categories_contributed is required for operation type: %s", op.OperationType)
			}
			if op.PrivacyParams.MinValue == nil {
				return fmt.Errorf("min_value is required for operation type: %s", op.OperationType)
			}
			if op.PrivacyParams.MaxValue == nil {
				return fmt.Errorf("max_value is required for operation type: %s", op.OperationType)
			}
		}
	}
	return nil
}

func LoadYamlConfig(filename string) (config *YamlConfig, err error) {
	log.Infof("Loading config from file: %s", filename)
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator(enLocale.Locale())
	err = registerCustomErrorMessages(validate, trans)
	if err != nil {
		return nil, fmt.Errorf("failed to register custom translations: %v", err)
	}

	dec := yaml.NewDecoder(file, yaml.Validator(validate), yaml.Strict())

	config = &YamlConfig{}
	err = dec.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	err = validateFieldsOpType(config)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	return config, nil
}
