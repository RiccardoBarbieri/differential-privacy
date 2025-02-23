package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/core/typex"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	log "github.com/golang/glog"
	"godp/model"
	"io"
	"os"
	"reflect"
	"strings"
)

func init() {
	register.Function2x0[string, beam.V](printKVConsoleFn)
	register.Function1x0[model.Admission](printAdmissionConsoleFn)
	register.Function1x0[model.ValuesStruct](printInterfaceConsoleFn)
	register.Function1x0[string](printStringConsoleFn)
}

func LoadCleanDataset(scope beam.Scope, fileName string) beam.PCollection {
	scope = scope.Scope("LoadCleanDataset")
	admissions := ReadInput(scope, fileName)
	admissionsCleaned := CleanDataset(scope, admissions)
	return admissionsCleaned
}

// ReadInput reads from input csv file and returns a beam PCollection of Admission
func ReadInput(scope beam.Scope, fileName string) beam.PCollection {
	scope = scope.Scope("readInput")
	lines := textio.Read(scope, fileName)
	return beam.ParDo(scope, model.CreateAdmissionFn, lines)
}

func ReadGenericInput(scope beam.Scope, fileName string) beam.PCollection {
	scope = scope.Scope("readGenericInput")
	lines := textio.Read(scope, fileName)
	return beam.ParDo(scope, model.CreateGenericStruct, lines)
}

func WriteOutput(scope beam.Scope, col beam.PCollection, fileName string) {
	scope = scope.Scope("WriteOutput")
	if typex.IsKV(col.Type()) {
		textio.Write(scope, fileName, beam.ParDo(scope, formatKVCsvFn, col))
	} else if col.Type().Type() == reflect.TypeOf(model.Admission{}) {
		textio.Write(scope, fileName, beam.ParDo(scope, formatStructCsvFn, col))
	} else {
		panic("unsupported output type: " + fmt.Sprintf("%T", col.Type()) + " filename: " + fileName)
	}
}

func WriteHeaders(fileName string, headers ...string) {
	if len(headers) > 0 {
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return
		}
		defer closeFile(file)
		reader := bufio.NewReader(file)
		// read lines until EOF
		var lines []string
		for i := 0; ; i++ {
			line, err := reader.ReadString('\n')
			lines = AppendStringArray(lines, line)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error reading from file: %v", err)
			}
			if strings.TrimSpace(line) == "" {
				break
			}
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			log.Fatalf("Error seeking to start of file: %v", err)
		}
		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(strings.Join(headers, ",") + "\n")
		if err != nil {
			log.Fatalf("Error writing to file: %v", err)
		}

		for _, line := range lines {
			_, err = writer.WriteString(line)
			if err != nil {
				log.Fatalf("Error writing to file: %v", err)
			}
		}
		err = writer.Flush()
		if err != nil {
			log.Fatalf("Error flushing file: %v", err)
		}

	}
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Warningf("Error closing file: %v", err)
	}
}

// PrintConsole prints the key-value pairs to the console.
func PrintConsole(scope beam.Scope, col beam.PCollection) {
	scope = scope.Scope("PrintConsole")

	if typex.IsKV(col.Type()) {
		beam.ParDo0(scope, printKVConsoleFn, col)
	} else if col.Type().Type() == reflect.TypeOf("") {
		beam.ParDo0(scope, printStringConsoleFn, col)
	} else if col.Type().Type() == reflect.TypeOf(model.Admission{}) {
		beam.ParDo0(scope, printAdmissionConsoleFn, col)
	} else if col.Type().Type() == reflect.TypeOf(model.ValuesStruct{}) {
		beam.ParDo0(scope, printInterfaceConsoleFn, col)
	}
}

// printKVConsoleFn is a DoFn that prints key-value pairs to the console and returns the printed value.
func printKVConsoleFn(k string, v beam.V) {
	fmt.Printf("%s -> %d\n", k, v)
}

// printAdmissionConsoleFn is a DoFn that prints Admission objects to the console and returns the printed value.
func printAdmissionConsoleFn(admission model.Admission) {
	fmt.Printf("%s\n", admission)
}

func printInterfaceConsoleFn(baseStruct model.ValuesStruct) {
	//fmt.Printf("%s\n", reflect.TypeOf(baseStruct.values["DateofAdmission"]))
	fmt.Printf("%v\n", baseStruct)
}

func printStringConsoleFn(s string) {
	fmt.Printf("%s\n", s)
}

func RemoveHeadersAndSaveCsv(filename string) (newFilename string, err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer closeFile(file)
	reader := bufio.NewReader(file)
	var lines []string
	for i := 0; ; i++ {
		line, err := reader.ReadString('\n')
		lines = AppendStringArray(lines, line)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(line) == "" {
			break
		}
	}
	err = file.Close()
	if err != nil {
		return "", err
	}

	newFilename = insertSuffixFilename(filename, "csv", "_noheader")
	_ = os.Remove(newFilename)
	newFile, err := os.OpenFile(newFilename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer closeFile(newFile)
	writer := bufio.NewWriter(newFile)
	for _, line := range lines[1:] {
		_, err = writer.WriteString(line)
		if err != nil {
			return "", err
		}
	}
	err = writer.Flush()
	if err != nil {
		return "", err
	}
	err = newFile.Close()
	if err != nil {
		return "", err
	}
	return newFilename, nil
}

func insertSuffixFilename(s string, ext string, suffix string) string {
	return strings.TrimSuffix(s, "."+ext) + suffix + "." + ext
}

// GetHeaders reads the first line of a CSV file and returns the headers.
//
// Parameters:
//   - filename: The path to the CSV file to read.
//
// Returns:
//   - []string: A slice containing the headers from the first line of the CSV file.
//   - error: An error if the file cannot be read or parsed, or nil if successful.
func GetHeaders(filename string) ([]string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer closeFile(file)
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(strings.NewReader(line))
	cols, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	var cleanCols []string
	for _, col := range cols {
		cleanCols = append(cleanCols, col)
	}
	return cleanCols, nil
}
