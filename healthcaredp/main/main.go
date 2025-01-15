package main

import (
	"fmt"
	"healthcaredp/command"
)

func main() {
	var err error

	err = command.RootCmd.Execute()
	if err != nil {
		_ = fmt.Errorf("error executing command: %v", err)
	}
}
