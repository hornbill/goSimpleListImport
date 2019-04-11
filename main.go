package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"

	_ "github.com/alexbrainman/odbc"     //ODBC Driver
	_ "github.com/denisenkom/go-mssqldb" //Microsoft SQL Server driver - v2005+
	_ "github.com/go-sql-driver/mysql"   //MySQL v4.1 to v5.x and MariaDB driver
	_ "github.com/hornbill/mysql320"     //MySQL v3.2.0 to v5 driver - Provides SWSQL (MySQL 4.0.16) support
)

// main package
func main() {
	//-- Start Time for Durration
	startTime = time.Now()
	localLogFileName = "list_import_" + time.Now().Format("20060102150405") + ".log"

	parseFlags()
	//-- If configVersion just output version number and die
	if configVersion {
		fmt.Printf("%v \n", version)
		return
	}
	//-- Output to CLI and Log
	logger(1, "---- Hornbill Simple List Import Utility V"+fmt.Sprintf("%v", version)+" ----", true)
	logger(1, "Flag - Config File "+configFileName, true)
	logger(1, "Flag - Dry Run "+fmt.Sprintf("%v", configDryRun), true)
	logger(1, "Flag - Concurrent Requests "+fmt.Sprintf("%v", configMaxRoutines), true)

	//-- Load Configuration File Into Struct
	importConf, boolConfLoaded = loadConfig()
	if !boolConfLoaded {
		logger(4, "Unable to load config, process closing.", true)
		return
	}

	configMaxRoutines = len(importConf.HBConf.APIKeys)
	if configMaxRoutines < 1 || configMaxRoutines > 10 {
		color.Red("The maximum allowed workers is between 1 and 10 (inclusive).\n\n")
		color.Red("You have included " + strconv.Itoa(configMaxRoutines) + " API keys. Please try again, with a valid number of keys.")
		return
	}

	//Set SQL driver ID string for Application Data
	if importConf.AppDBConf.Driver == "" {
		logger(4, "AppDBConf SQL Driver not set in configuration.", true)
		return
	}
	if importConf.AppDBConf.Driver == "swsql" {
		appDBDriver = "mysql320"
	} else if importConf.AppDBConf.Driver == "mysql" ||
		importConf.AppDBConf.Driver == "mssql" ||
		importConf.AppDBConf.Driver == "mysql320" ||
		importConf.AppDBConf.Driver == "odbc" ||
		importConf.AppDBConf.Driver == "csv" {
		appDBDriver = importConf.AppDBConf.Driver
	} else {
		logger(4, "The driver ("+importConf.AppDBConf.Driver+") for the Application Database specified in the configuration file is not valid.", true)
		return
	}

	//-- Build DB connection string
	connStrAppDB = buildConnectionString()

	//Get request type import config, process each in turn
	for _, val := range importConf.Lists {
		if val.ListName != "" && val.Application != "" {

			if val.Rebuild {
				// delete all items from list
				emptySimpleList(val.Application, val.ListName)
			}
			
			if appDBDriver == "odbc" ||
				appDBDriver == "xls" ||
				appDBDriver == "csv" {
				processSimpleListDataODBC(val)
			} else {
				processSimpleListData(val)

			}

		} else {
			// skip whatever
			continue
		}
	}

	//-- End output
	logger(3, "", true)
	logger(3, "Lists Deleted: "+fmt.Sprintf("%d", counters.deleted), true)
	logger(3, "Lists not deleted: "+fmt.Sprintf("%d", counters.deletedSkipped), true)
	logger(3, "Items Created: "+fmt.Sprintf("%d", counters.created), true)
	logger(3, "Items Skipped: "+fmt.Sprintf("%d", counters.createdSkipped), true)
	//-- Show Time Takens
	endTime = time.Since(startTime)
	logger(3, "Time Taken: "+fmt.Sprintf("%v", endTime), true)
	logger(3, "---- Hornbill Simple List Import Complete ---- ", true)
}
