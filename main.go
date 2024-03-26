// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"log"
	"os"

	"github.com/vmware-labs/marketplace-cli/v2/cmd"
)

func main() {
	// Open a log file for writing
	log.Println("Creating the file logfile and printing the lgos..")
	logFile, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	// Set the log output to the log file
	log.SetOutput(logFile)

	// start the cli
	cmd.Execute()
}
