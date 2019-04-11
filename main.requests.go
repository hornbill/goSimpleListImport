package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
//	"strconv"
//	"strings"
	"sync"
//	"time"

	apiLib "github.com/hornbill/goApiLib"
//	"github.com/hornbill/pb"
)


//logNewCallJobs - Function takes external call data in a map, and logs to Hornbill
func logNewCallJobs(jobs chan RequestDetails, wg *sync.WaitGroup, espXmlmc *apiLib.XmlmcInstStruct, ltpDetails *listToProcessStruct) {
	defer wg.Done()
	for request := range jobs {
		var buffer bytes.Buffer
		buffer.WriteString(loggerGen(3, "   "))
		
		_ = addItem(request, espXmlmc, &buffer, ltpDetails)

		bufferMutex.Lock()
		loggerWriteBuffer(buffer.String())
		bufferMutex.Unlock()
		buffer.Reset()
	}
}

func addItem(request RequestDetails, espXmlmc *apiLib.XmlmcInstStruct, buffer *bytes.Buffer, ltpDetails *listToProcessStruct) (bool) {

	//Loop through core fields from config, add to XMLMC Params
/*	for k, v := range mapGenericConf.CoreFieldMapping {
		boolAutoProcess := true
		strAttribute = fmt.Sprintf("%v", k)
		strMapping = fmt.Sprintf("%v", v)

	}
*/
	strVal := getFieldValue(ltpDetails.ItemValue, &request.CallMap)
	strDisp := getFieldValue(ltpDetails.DefaultDisplay, &request.CallMap)
	if strVal != "" && strDisp != "" {
		espXmlmc.SetParam("application", "com.hornbill.servicemanager")
		espXmlmc.SetParam("listName", ltpDetails.ListName)
		espXmlmc.SetParam("itemValue", strVal)
		espXmlmc.SetParam("defaultItemName", strDisp)
	} else {
	
		mutexCounters.Lock()
		counters.createdSkipped++
		mutexCounters.Unlock()
		buffer.WriteString(loggerGen(4, "Skipping Missing DATA : " + strVal + " - " + strDisp))
		return false
	
	}
	for _, oTranslations := range ltpDetails.Translations {
		espXmlmc.OpenElement("itemNameTranslation")
			espXmlmc.SetParam("itemName",  getFieldValue(oTranslations.Display, &request.CallMap))
			espXmlmc.SetParam("language",  getFieldValue(oTranslations.Language, &request.CallMap))
		espXmlmc.CloseElement("itemNameTranslation")
	}


	//-- Check for Dry Run
	if !configDryRun {
		XMLRequest := espXmlmc.GetParam()
		XMLCreate, xmlmcErr := espXmlmc.Invoke("data", "listAddItem")
		if xmlmcErr != nil {

			mutexCounters.Lock()
			counters.createdSkipped++
			mutexCounters.Unlock()
			buffer.WriteString(loggerGen(4, "API Call Failed: New Item : "+xmlmcErr.Error()))
			buffer.WriteString(loggerGen(1, "[XML] "+XMLRequest))
			return false
		}
		var xmlRespon xmlmcRequestResponseStruct

		err := xml.Unmarshal([]byte(XMLCreate), &xmlRespon)
		if err != nil {
			mutexCounters.Lock()
			counters.createdSkipped++
			mutexCounters.Unlock()
			buffer.WriteString(loggerGen(4, "Response Unmarshal failed: New Item : "+fmt.Sprintf("%v", err)))
			buffer.WriteString(loggerGen(1, "[XML] "+XMLRequest))
			return false
		}
		if xmlRespon.MethodResult != "ok" {
			mutexCounters.Lock()
			counters.createdSkipped++
			mutexCounters.Unlock()
			buffer.WriteString(loggerGen(4, "MethodResult not OK: New Item : "+xmlRespon.State.ErrorRet))
			buffer.WriteString(loggerGen(1, "[XML] "+XMLRequest))
			return false
		}
		
		mutexCounters.Lock()
		counters.created++
		mutexCounters.Unlock()

	} else {
		//-- DEBUG XML TO LOG FILE
		var XMLSTRING = espXmlmc.GetParam()
		buffer.WriteString(loggerGen(1, "Request Log XML "+XMLSTRING))
		mutexCounters.Lock()
		counters.createdSkipped++
		mutexCounters.Unlock()
		espXmlmc.ClearParam()
	}
	return true
}

