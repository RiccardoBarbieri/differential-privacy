package main

import (
	"github.com/golang/glog"
	"io"
	"log"
	"os"

	"godp/command"
)

func main() {
	//file, err := os.OpenFile("/var/log/"+command.RootCmd.Name()+"/"+command.RootCmd.Name()+".log", os.O_CREATE, 0644)
	//if err != nil {
	//	log.Fatalf("Error opening log file: %v", err)
	//}
	//defer func() {
	//	err := file.Close()
	//	if err != nil {
	//		return
	//	}
	//}()

	log.SetOutput(io.Discard)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//log.SetOutput(file)

	//flag.Set("logtostderr", "true")
	//flag.Set("stderrthreshold", "INFO")

	err := command.RootCmd.Execute()
	if err != nil {
		glog.Errorf("error executing command: %v", err)
		os.Exit(1)
	}

	glog.Flush()
}
