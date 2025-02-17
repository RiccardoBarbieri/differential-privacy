package model

import (
	"encoding/csv"
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	"github.com/ompluscator/dynamic-struct"
	"strconv"
	"strings"
	"time"
)

func init() {
	register.Function2x1[string, func(struc ValuesStruct), error](CreateGenericStruct)
	register.Emitter1[ValuesStruct]()
}

var GenericStruct dynamicstruct.DynamicStruct
var Headers []string
var TypesMap map[string]string
var IdFieldIndex int

func CompileTypesMap(types []TypeType) (map[string]string, error) {
	typesMap := make(map[string]string)
	for _, typ := range types {
		typesMap[typ.Column] = typ.Type
	}
	return typesMap, nil
}

func formatValue(valueStr string, typeStr string) (any, error) {
	switch typeStr {
	case "int":
		return strconv.Atoi(valueStr)
	case "string":
		return valueStr, nil
	case "bool":
		return strconv.ParseBool(valueStr)
	case "float":
		return strconv.ParseFloat(valueStr, 64)
	case "date":
		return time.Parse(time.DateOnly, valueStr)
	case "time":
		return time.Parse(time.TimeOnly, valueStr)
	case "datetime":
		return time.Parse(time.DateTime, valueStr)
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

type ValuesStruct struct {
	Values map[string]string
	Id     string
}

func CreateGenericStruct(line string, emit func(struc ValuesStruct)) error {
	reader := csv.NewReader(strings.NewReader(line))
	cols, err := reader.Read()
	if err != nil {
		return err
	}
	if len(cols) != len(Headers) {
		return fmt.Errorf("line containse %d columns, struct expects %d columns - %s", len(cols), len(Headers), line)
	}
	baseStruct := ValuesStruct{}
	baseStruct.Values = make(map[string]string)
	for i, col := range cols {
		//if typ, ok := TypesMap[Headers[i]]; ok {
		//	value, err := formatValue(col, typ)
		//	if err != nil {
		//		return err
		//	}
		//	baseStruct.Values[Headers[i]] = value
		//} else {
		//	baseStruct.Values[Headers[i]] = col
		//}
		baseStruct.Values[Headers[i]] = col
		baseStruct.Id = cols[IdFieldIndex]
	}

	emit(baseStruct)
	return nil
}
