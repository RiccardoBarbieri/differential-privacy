package main

import (
	"fmt"
	"healthcaredp/cmds"
)

func main() {
	var err error

	err = cmds.RootCmd.Execute()
	if err != nil {
		_ = fmt.Errorf("error executing command: %v", err)
	}
}
