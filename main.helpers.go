package main

import (
//	"bytes"
	"encoding/json"
//	"encoding/xml"
	"flag"
	"fmt"
	"os"
//	"strconv"
	"strings"
//	"time"

	hornbillHelpers "github.com/hornbill/goHornbillHelpers"
)

//loadConfig -- Function to Load Configruation File
func loadConfig() (importConfStruct, bool) {
	boolLoadConf := true
	//-- Check Config File File Exists
	cwd, _ := os.Getwd()
	configurationFilePath := cwd + "/" + configFileName
	logger(1, "Loading Config File: "+configurationFilePath, false)
	if _, fileCheckErr := os.Stat(configurationFilePath); os.IsNotExist(fileCheckErr) {
		logger(4, "No Configuration File", true)
		os.Exit(102)
	}
	//-- Load Config File
	file, fileError := os.Open(configurationFilePath)
	//-- Check For Error Reading File
	if fileError != nil {
		logger(4, "Error Opening Configuration File: "+fmt.Sprintf("%v", fileError), true)
		boolLoadConf = false
	}

	//-- New Decoder
	decoder := json.NewDecoder(file)
	//-- New Var based on importConfStruct
	edbConf := importConfStruct{}
	//-- Decode JSON
	err := decoder.Decode(&edbConf)
	//-- Error Checking
	if err != nil {
		logger(4, "Error Decoding Configuration File: "+fmt.Sprintf("%v", err), true)
		boolLoadConf = false
	}
	//-- Return New Config
	return edbConf, boolLoadConf
}


//parseFlags - grabs and parses command line flags
func parseFlags() {
	flag.StringVar(&configFileName, "file", "conf.json", "Name of the configuration file to load")
	flag.BoolVar(&configDryRun, "dryrun", false, "Dump import XML to log instead of creating requests")
	flag.BoolVar(&configVersion, "version", false, "Outputs version number and exits tool")
	flag.Parse()
}


func logger(t int, s string, outputToCLI bool) {
	hornbillHelpers.Logger(t, s, outputToCLI, localLogFileName)
}

func loggerGen(t int, s string) string {

	var errorLogPrefix = ""
	//-- Create Log Entry
	switch t {
	case 1:
		errorLogPrefix = "[DEBUG] "
	case 2:
		errorLogPrefix = "[MESSAGE] "
	case 3:
		errorLogPrefix = ""
	case 4:
		errorLogPrefix = "[ERROR] "
	case 5:
		errorLogPrefix = "[WARNING] "
	case 6:
		errorLogPrefix = ""
	}
	return errorLogPrefix + s + "\n\r"
}
func loggerWriteBuffer(s string) {
	if s != "" {
		logLines := strings.Split(s, "\n\r")
		for _, line := range logLines {
			if line != "" {
				logger(0, line, false)
			}
		}
	}
}
