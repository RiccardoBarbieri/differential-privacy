package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
)

func init() {
	register.DoFn2x1[ValuesStruct, func(ValuesStruct), error](&FilterFn{})
	register.Emitter1[ValuesStruct]()
}

// FilterFn is a DoFn that filters records based on configured filters
type FilterFn struct {
	Filters  []FilterType
	TypesMap map[string]string
}

// ProcessElement evaluates all filters on a record
func (fn *FilterFn) ProcessElement(record ValuesStruct, emit func(ValuesStruct)) error {
	// If no filters, emit the record
	if len(fn.Filters) == 0 {
		emit(record)
		return nil
	}

	// Check all filters - all must pass (AND logic)
	for _, filter := range fn.Filters {
		value, exists := record.Values[filter.Column]
		if !exists {
			// Column doesn't exist, filter fails
			return nil
		}

		pass, err := evaluateFilter(value, filter, fn.TypesMap)
		if err != nil {
			return err
		}
		if !pass {
			return nil
		}
	}

	emit(record)
	return nil
}

// evaluateFilter evaluates a single filter condition
func evaluateFilter(value string, filter FilterType, typesMap map[string]string) (bool, error) {
	colType, ok := typesMap[filter.Column]
	if !ok {
		// Default to string comparison if type not specified
		colType = "string"
	}

	switch colType {
	case "int", "float":
		return evaluateNumericFilter(value, filter)
	case "string":
		return evaluateStringFilter(value, filter)
	default:
		// For other types (date, time, etc.), treat as string comparison
		return evaluateStringFilter(value, filter)
	}
}

// evaluateNumericFilter evaluates numeric comparisons
func evaluateNumericFilter(value string, filter FilterType) (bool, error) {
	// Parse the record value
	recVal, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return false, fmt.Errorf("failed to parse numeric value '%s' for column '%s': %v", value, filter.Column, err)
	}

	// Parse the filter value
	filterVal, err := strconv.ParseFloat(strings.TrimSpace(filter.Value), 64)
	if err != nil {
		return false, fmt.Errorf("failed to parse filter value '%s' for column '%s': %v", filter.Value, filter.Column, err)
	}

	switch filter.Operator {
	case "=":
		return recVal == filterVal, nil
	case "!=":
		return recVal != filterVal, nil
	case "<":
		return recVal < filterVal, nil
	case ">":
		return recVal > filterVal, nil
	case "<=":
		return recVal <= filterVal, nil
	case ">=":
		return recVal >= filterVal, nil
	default:
		return false, fmt.Errorf("unsupported operator '%s' for numeric comparison", filter.Operator)
	}
}

// evaluateStringFilter evaluates string comparisons
func evaluateStringFilter(value string, filter FilterType) (bool, error) {
	recVal := strings.TrimSpace(value)
	filterVal := strings.TrimSpace(filter.Value)

	switch filter.Operator {
	case "=":
		return recVal == filterVal, nil
	case "!=":
		return recVal != filterVal, nil
	case "<":
		return recVal < filterVal, nil
	case ">":
		return recVal > filterVal, nil
	case "<=":
		return recVal <= filterVal, nil
	case ">=":
		return recVal >= filterVal, nil
	default:
		return false, fmt.Errorf("unsupported operator '%s' for string comparison", filter.Operator)
	}
}

// ApplyFilters applies all configured filters to a PCollection
func ApplyFilters(scope beam.Scope, col beam.PCollection, filters []FilterType, typesMap map[string]string) beam.PCollection {
	if len(filters) == 0 {
		return col
	}
	scope = scope.Scope("ApplyFilters")
	return beam.ParDo(scope, &FilterFn{Filters: filters, TypesMap: typesMap}, col)
}

// ValidateFilterColumns checks that all filter columns are specified in the types map
func ValidateFilterColumns(filters []FilterType, types []TypeType) error {
	typesMap := make(map[string]bool)
	for _, t := range types {
		typesMap[t.Column] = true
	}

	for _, filter := range filters {
		if !typesMap[filter.Column] {
			return fmt.Errorf("filter column '%s' must be specified in the types section", filter.Column)
		}
	}
	return nil
}
