package main

import (
//	"bytes"
//	"encoding/xml"
//	"fmt"
//	"strconv"

	"github.com/hornbill/goApiLib"
)

//NewEspXmlmcSession - New Xmlmc Session variable (Cloned Session)
func NewEspXmlmcSession(apiKey string) *apiLib.XmlmcInstStruct {
	espXmlmcLocal := apiLib.NewXmlmcInstance(importConf.HBConf.InstanceID)
	espXmlmcLocal.SetAPIKey(apiKey)
	return espXmlmcLocal
}
