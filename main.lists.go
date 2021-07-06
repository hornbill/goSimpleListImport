package main

import (
	"bytes"
	"encoding/xml"
	"fmt"

	//	"strconv"
	//	"strings"
	//	"sync"
	//	"time"
	"github.com/hornbill/pb"
	//	apiLib "github.com/hornbill/goApiLib"
)

func emptySimpleList(applicationName string, listName string) {
	espXmlmc := NewEspXmlmcSession(importConf.HBConf.APIKeys[0])
	espXmlmc.SetParam("application", applicationName)
	espXmlmc.SetParam("listName", listName)
	if !configDryRun {
		XMLRequest := espXmlmc.GetParam()
		XMLCreate, xmlmcErr := espXmlmc.Invoke("data", "listDelete")
		if xmlmcErr != nil {
			mutexCounters.Lock()
			counters.deletedSkipped++
			mutexCounters.Unlock()
			//			buffer.WriteString(loggerGen(4, "API Call Failed: Delete List : "+xmlmcErr.Error()))
			//			buffer.WriteString(loggerGen(1, "[XML] "+XMLRequest))
			_ = loggerGen(4, "API Call Failed: Delete List : "+xmlmcErr.Error())
			_ = loggerGen(1, "[XML] "+XMLRequest)
			return
		}
		var xmlRespon xmlmcRequestResponseStruct
		err := xml.Unmarshal([]byte(XMLCreate), &xmlRespon)
		if err != nil {
			mutexCounters.Lock()
			counters.deletedSkipped++
			mutexCounters.Unlock()
			//buffer.WriteString(loggerGen(4, "Response Unmarshal failed: Deleted List : "+fmt.Sprintf("%v", err)))
			//buffer.WriteString(loggerGen(1, "[XML] "+XMLRequest))
			_ = loggerGen(4, "Response Unmarshal failed: Deleted List : "+fmt.Sprintf("%v", err))
			_ = loggerGen(1, "[XML] "+XMLRequest)
			return
		}
		if xmlRespon.MethodResult != "ok" {
			mutexCounters.Lock()
			counters.deletedSkipped++
			mutexCounters.Unlock()
			//buffer.WriteString(loggerGen(4, "MethodResult not OK: Deleted List : "+xmlRespon.State.ErrorRet))
			//buffer.WriteString(loggerGen(1, "[XML] "+XMLRequest))
			_ = loggerGen(4, "MethodResult not OK: Deleted List : "+xmlRespon.State.ErrorRet)
			_ = loggerGen(1, "[XML] "+XMLRequest)
			return
		}
		mutexCounters.Lock()
		counters.deleted++
		mutexCounters.Unlock()
	} else {
		//-- DEBUG XML TO LOG FILE
		var XMLSTRING = espXmlmc.GetParam()
		//buffer.WriteString(loggerGen(1, "Delete Log XML "+XMLSTRING))
		_ = loggerGen(1, "Delete Log XML "+XMLSTRING)
		mutexCounters.Lock()
		counters.deletedSkipped++
		mutexCounters.Unlock()
		espXmlmc.ClearParam()
	}
}

//processCallData - Query External call data, process accordingly
func processSimpleListData(ltpDetails listToProcessStruct) {
	arrCallDetailsMaps, success := queryDBListDetails(ltpDetails.SQL)
	if success {

		bar := pb.StartNew(len(arrCallDetailsMaps))
		defer bar.FinishPrint(ltpDetails.Application + ":" + ltpDetails.ListName + " List Import Complete")

		jobs := make(chan RequestDetails, configMaxRoutines)

		for w := 0; w < configMaxRoutines; w++ {
			wg.Add(1)
			espXmlmc := NewEspXmlmcSession(importConf.HBConf.APIKeys[w])
			go logNewCallJobs(jobs, &wg, espXmlmc, &ltpDetails)
		}

		for _, callRecord := range arrCallDetailsMaps {
			mutexBar.Lock()
			bar.Increment()
			mutexBar.Unlock()
			jobs <- RequestDetails{CallMap: callRecord}
		}

		close(jobs)
		wg.Wait()

	} else {
		logger(4, "Request search failed for type: "+ltpDetails.Application+":"+ltpDetails.ListName, true)
	}
}

//processCallDataODBC - Query ODBC call data, process accordingly
func processSimpleListDataODBC(ltpDetails listToProcessStruct) {
	arrCallDetailsMaps, success := queryDBListDetails(ltpDetails.SQL)
	if success {
		bar := pb.StartNew(len(arrCallDetailsMaps))
		defer bar.FinishPrint(ltpDetails.Application + ":" + ltpDetails.ListName + " List Import Complete")
		espXmlmc := NewEspXmlmcSession(importConf.HBConf.APIKeys[0])

		for _, callRecord := range arrCallDetailsMaps {
			mutexBar.Lock()
			bar.Increment()
			mutexBar.Unlock()

			var buffer bytes.Buffer

			//fmt.Println("%v", callRecord)
			//fmt.Println("%v", espXmlmc)

			buffer.WriteString(loggerGen(3, "   "))
			_ = addItem(RequestDetails{CallMap: callRecord}, espXmlmc, &buffer, &ltpDetails)
			//oldCallRef, newCallRef, oldCallGUID = logNewCall(RequestDetails{GenericImportConf: mapGenericConf, CallMap: callRecord}, espXmlmc, &buffer)

			bufferMutex.Lock()
			loggerWriteBuffer(buffer.String())
			bufferMutex.Unlock()
			//buffer.Reset()
		}
	} else {
		logger(4, "Request search failed for type: "+ltpDetails.Application+":"+ltpDetails.ListName, true)
	}
}
