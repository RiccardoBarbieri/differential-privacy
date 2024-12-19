package utils

import (
	"bufio"
	"fmt"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/core/typex"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	log "github.com/golang/glog"
	"healthcaredp"

	"io"
	"os"
	"reflect"
	"strings"
)

func init() {

	register.Function2x0[string, beam.V](printKVConsoleFn)
	register.Function1x0[healtcaredp.Admission](printAdmissionConsoleFn)
}

// Reads from input csv file and returns a beam PCollection of Admission
func ReadInput(scope beam.Scope, fileName string) beam.PCollection {
	scope = scope.Scope("readInput")
	lines := textio.Read(scope, fileName)
	return beam.ParDo(scope, healtcaredp.CreateAdmissionFn, lines)
}

func WriteOutput(scope beam.Scope, col beam.PCollection, fileName string, headers ...string) {
	scope = scope.Scope("writeOutput")
	if typex.IsKV(col.Type()) {
		textio.Write(scope, fileName, beam.ParDo(scope, formatKVCsvFn, col))

	} else if col.Type().Type() == reflect.TypeOf(healtcaredp.Admission{}) {
		textio.Write(scope, fileName, beam.ParDo(scope, formatStructCsvFn, col))
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
		log.Fatalf("Error closing file: %v", err)
	}
}

// PrintConsole prints the key-value pairs to the console.
func PrintConsole(scope beam.Scope, col beam.PCollection) {
	scope = scope.Scope("PrintConsole")
	if typex.IsKV(col.Type()) {
		beam.ParDo0(scope, printKVConsoleFn, col)
	} else if col.Type().Type() == reflect.TypeOf(healtcaredp.Admission{}) {
		beam.ParDo0(scope, printAdmissionConsoleFn, col)
	}
}

// printKVConsoleFn is a DoFn that prints key-value pairs to the console and returns the printed value.
func printKVConsoleFn(k string, v beam.V) {
	fmt.Printf("%s -> %d\n", k, v)
}

// printAdmissionConsoleFn is a DoFn that prints Admission objects to the console and returns the printed value.
func printAdmissionConsoleFn(admission healtcaredp.Admission) {
	fmt.Printf("%s\n", admission)
}
