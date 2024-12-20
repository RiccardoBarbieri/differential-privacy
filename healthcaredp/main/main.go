package main

import (
	"fmt"
	"healthcaredp/cmd"
)

func main() {
	var err error

	err = cmd.RootCmd.Execute()
	if err != nil {
		_ = fmt.Errorf("error executing command: %v", err)
	}
}
