package utils

import (
	"encoding/csv"
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	log "github.com/golang/glog"

	"reflect"
	"strconv"
	"strings"
	"time"

	// required for local file access
	_ "github.com/apache/beam/sdks/v2/go/pkg/beam/io/filesystem/local"
)

func init() {
	register.Function2x1[string, beam.V, string](formatKVCsvFn)
	register.Function1x1[interface{}, string](formatStructCsvFn)
}

func formatKVCsvFn(k string, v beam.V) string {
	return fmt.Sprintf("%s,%d", k, v)
}

// formatStructCsvFn converts a struct to a CSV-formatted string.
//
// This function takes an interface{} parameter, which is expected to be a struct,
// and converts it to a CSV-formatted string. Each field of the struct is
// converted to a string representation and separated by commas.
//
// Parameters:
//   - s: An interface{} that should be a struct. If it's not a struct,
//     the function will log a fatal error and exit.
//
// Returns:
//
//	A string containing the CSV representation of the input struct.
//	Each field of the struct is converted to a string and separated by commas.
//
// Note: This function prints the resulting CSV string to stdout before returning it.
func formatStructCsvFn(s interface{}) string {
	reflectValue := reflect.ValueOf(s)
	if !(reflectValue.Kind() == reflect.Struct) {
		log.Exitf("s must be a struct, got %T", s)
	}
	sb := strings.Builder{}
	writer := csv.NewWriter(&sb)
	var fields = make([]string, reflectValue.NumField())
	for i := 0; i < reflectValue.NumField(); i++ {
		fields[i] = formatType(reflectValue.Field(i))
	}
	err := writer.Write(fields)
	if err != nil {
		log.Fatalf("Error writing headers to CSV: %v", err)
	}
	writer.Flush()
	return sb.String()[:sb.Len()-1]
}

// formatType converts a reflect.Value to its string representation.
//
// This function handles various types including strings, integers, floats, and time.Time.
// For unsupported types, it returns "TYPE_NOT_IMPLEMENTED".
//
// Parameters:
//   - t: A reflect.Value representing the value to be formatted.
//
// Returns:
//
//	A string representation of the input value. For unsupported types, it returns "TYPE_NOT_IMPLEMENTED".
func formatType(t reflect.Value) string {
	switch t.Kind() {
	case reflect.String:
		return t.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(t.Int(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(t.Float(), 'f', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(t.Float(), 'f', -1, 32)
	case reflect.Struct:
		if t.Type() == reflect.TypeOf(time.Time{}) {
			return t.Interface().(time.Time).Format(time.DateOnly)
		}
	default:
		return "TYPE_NOT_IMPLEMENTED"
	}
	return "TYPE_NOT_IMPLEMENTED"
}

func StructCsvHeaders(s interface{}) []string {
	reflectValue := reflect.ValueOf(s)
	if !(reflectValue.Kind() == reflect.Struct) {
		log.Exitf("s must be a struct, got %T", s)
	}
	fields := make([]string, 0, reflectValue.NumField())
	for i := 0; i < reflectValue.NumField(); i++ {
		fields = append(fields, reflectValue.Type().Field(i).Name)
	}
	return fields
}

func AppendStringArray(a []string, s string) []string {
	if len(a)+1 > cap(a) {
		newArray := make([]string, 2*len(a))
		copy(newArray, a)
		a = newArray
	}
	a = append(a, s)
	return a
}
