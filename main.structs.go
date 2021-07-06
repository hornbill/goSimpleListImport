package main

import (
	"sync"
	"time"
)

const (
	version = "1.0.2"
)

var (
	localLogFileName  string
	appDBDriver       string
	boolConfLoaded    bool
	configFileName    string
	configDryRun      bool
	configMaxRoutines int
	configVersion     bool
	connStrAppDB      string
	counters          counterTypeStruct
	sqlCallQuery      string
	importConf        importConfStruct
	startTime         time.Time
	endTime           time.Duration
	wg                sync.WaitGroup
	bufferMutex       = &sync.Mutex{}
	mutexBar          = &sync.Mutex{}
	mutexCounters     = &sync.Mutex{}
)

// ----- Structures -----
type counterTypeStruct struct {
	sync.Mutex
	created        int
	createdSkipped int
	deleted        int
	deletedSkipped int
}

//----- Config Data Structs
type importConfStruct struct {
	HBConf    hbConfStruct    //Hornbill Instance connection details
	AppDBConf appDBConfStruct //App Data (swdata) connection details
	Lists     []listToProcessStruct
}
type listToProcessStruct struct {
	SQL            string
	Rebuild        bool
	Application    string
	ListName       string
	ItemValue      string
	DefaultDisplay string
	Translations   []translationsStruct
}

type translationsStruct struct {
	Language string
	Display  string
}

type hbConfStruct struct {
	InstanceID string
	APIKeys    []string
}

type appDBConfStruct struct {
	Address  string
	Driver   string
	Server   string
	UserName string
	Password string
	Port     int
	Database string
	Encrypt  bool
}

//----- Shared Structs -----
type stateStruct struct {
	Code     string `xml:"code"`
	ErrorRet string `xml:"error"`
}

//----- Data Structs -----

//----- Request Logged Structs
type xmlmcRequestResponseStruct struct {
	MethodResult     string      `xml:"status,attr"`
	RequestID        string      `xml:"params>primaryEntityData>record>h_pk_reference"`
	HistoricUpdateID string      `xml:"params>primaryEntityData>record>h_pk_updateid"`
	SiteCountry      string      `xml:"params>rowData>row>h_country"`
	Diags            []string    `xml:"diagnostic>log"`
	State            stateStruct `xml:"state"`
}

//RequestDetails struct for chan
type RequestDetails struct {
	CallMap map[string]interface{}
}
