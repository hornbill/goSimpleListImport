package main

import (
	"sync"
	"time"
)

const (
	version           = "1.0.1"
)

var (
	localLogFileName     string
	appDBDriver          string
	arrSpawnBPMs         []spawnBPMStruct
	boolConfLoaded       bool
	configFileName       string
	configDryRun         bool
	configMaxRoutines    int
	configVersion        bool
	connStrAppDB         string
	counters             counterTypeStruct
	sqlCallQuery         string
	importConf           importConfStruct
	timeNow              string
	startTime            time.Time
	endTime              time.Duration
	wg                   sync.WaitGroup
	bufferMutex          = &sync.Mutex{}
	mutexAnalysts        = &sync.Mutex{}
	mutexBar             = &sync.Mutex{}
	mutexCategories      = &sync.Mutex{}
	mutexCloseCategories = &sync.Mutex{}
	mutexCounters        = &sync.Mutex{}
	mutexCustomers       = &sync.Mutex{}
	mutexPriorities      = &sync.Mutex{}
	mutexServices        = &sync.Mutex{}
	mutexSites           = &sync.Mutex{}
	mutexTeams           = &sync.Mutex{}
	mutexRequests        = &sync.Mutex{}
	reqPrefix            string
)

// ----- Structures -----
type counterTypeStruct struct {
	sync.Mutex
	created         int
	createdSkipped  int
	deleted int
	deletedSkipped int
}

//----- Config Data Structs
type importConfStruct struct {
	HBConf                    hbConfStruct //Hornbill Instance connection details
	AppDBConf                 appDBConfStruct //App Data (swdata) connection details
	Lists []listToProcessStruct
}
type listToProcessStruct struct {
	SQL string
	Rebuild bool
	Application string
	ListName string
	ItemValue string
	DefaultDisplay string
	Translations []translationsStruct
}
 
type translationsStruct struct {
	Language string
	Display string
}

type hbConfStruct struct {
	InstanceID string
	APIKeys    []string
}

type spawnBPMStruct struct {
	RequestID string
	BPMID     string
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

type xmlmcResponse struct {
	MethodResult string      `xml:"status,attr"`
	State        stateStruct `xml:"state"`
}

//----- Shared Structs -----
type stateStruct struct {
	Code     string `xml:"code"`
	ErrorRet string `xml:"error"`
}

//----- Data Structs -----

type xmlmcSysSettingResponse struct {
	MethodResult string      `xml:"status,attr"`
	State        stateStruct `xml:"state"`
	Setting      string      `xml:"params>option>value"`
}

//----- Request Logged Structs
type xmlmcRequestResponseStruct struct {
	MethodResult     string      `xml:"status,attr"`
	RequestID        string      `xml:"params>primaryEntityData>record>h_pk_reference"`
	HistoricUpdateID string      `xml:"params>primaryEntityData>record>h_pk_updateid"`
	SiteCountry      string      `xml:"params>rowData>row>h_country"`
	Diags            []string    `xml:"diagnostic>log"`
	State            stateStruct `xml:"state"`
}
type xmlmcBPMSpawnedStruct struct {
	MethodResult string      `xml:"status,attr"`
	Identifier   string      `xml:"params>identifier"`
	State        stateStruct `xml:"state"`
}

//----- Site Structs
type siteListStruct struct {
	SiteName string
	SiteID   int
}
type xmlmcSiteListResponse struct {
	MethodResult string      `xml:"status,attr"`
	SiteID       int         `xml:"params>rowData>row>h_id"`
	SiteName     string      `xml:"params>rowData>row>h_site_name"`
	SiteCountry  string      `xml:"params>rowData>row>h_country"`
	State        stateStruct `xml:"state"`
}

//----- Priority Structs
type priorityListStruct struct {
	PriorityName string
	PriorityID   int
}
type xmlmcPriorityListResponse struct {
	MethodResult string      `xml:"status,attr"`
	PriorityID   int         `xml:"params>rowData>row>h_pk_priorityid"`
	PriorityName string      `xml:"params>rowData>row>h_priorityname"`
	State        stateStruct `xml:"state"`
}

//----- Service Structs
type serviceListStruct struct {
	ServiceName          string
	ServiceID            int
	ServiceBPMIncident   string
	ServiceBPMService    string
	ServiceBPMChange     string
	ServiceBPMProblem    string
	ServiceBPMKnownError string
	ServiceBPMRelease    string
	CatalogItems         []catalogItemListStruct
}

type catalogItemListStruct struct {
	CatalogItemName string `xml:"catalog_title"`
	CatalogItemID   int    `xml:"h_id"`
	RequestType     string `xml:"h_request_type"`
	BPM             string `xml:"h_bpm"`
	Status          string `xml:"h_status"`
}

type xmlmcServiceListResponse struct {
	MethodResult  string      `xml:"status,attr"`
	ServiceID     int         `xml:"params>rowData>row>h_linked_service_id"`
	ServiceName   string      `xml:"params>rowData>row>h_servicename"`
	BPMIncident   string      `xml:"params>rowData>row>h_incident_bpm_name"`
	BPMService    string      `xml:"params>rowData>row>h_service_bpm_name"`
	BPMChange     string      `xml:"params>rowData>row>h_change_bpm_name"`
	BPMProblem    string      `xml:"params>rowData>row>h_problem_bpm_name"`
	BPMKnownError string      `xml:"params>rowData>row>h_knownerror_bpm_name"`
	BPMRelease    string      `xml:"params>rowData>row>h_release_bpm_name"`
	State         stateStruct `xml:"state"`
}

type xmlmcCatalogItemListResponse struct {
	MethodResult string                  `xml:"status,attr"`
	CatalogItems []catalogItemListStruct `xml:"params>rowData>row"`
	FoundRows    int                     `xml:"params>foundRows"`
	State        stateStruct             `xml:"state"`
}

//----- Team Structs
type teamListStruct struct {
	TeamName string
	TeamID   string
}
type xmlmcTeamListResponse struct {
	MethodResult string      `xml:"status,attr"`
	TeamID       string      `xml:"params>rowData>row>h_id"`
	TeamName     string      `xml:"params>rowData>row>h_name"`
	State        stateStruct `xml:"state"`
}

//----- Category Structs
type categoryListStruct struct {
	CategoryCode string
	CategoryID   string
	CategoryName string
}
type xmlmcCategoryListResponse struct {
	MethodResult string      `xml:"status,attr"`
	CategoryID   string      `xml:"params>id"`
	CategoryName string      `xml:"params>fullname"`
	State        stateStruct `xml:"state"`
}



type xmlmcCustomerListResponse struct {
	MethodResult      string      `xml:"status,attr"`
	CustomerID        int         `xml:"params>rowData>row>h_pk_id"`
	CustomerFirstName string      `xml:"params>rowData>row>h_firstname"`
	CustomerLastName  string      `xml:"params>rowData>row>h_lastname"`
	State             stateStruct `xml:"state"`
}

//RequestDetails struct for chan
type RequestDetails struct {
	CallMap           map[string]interface{}
}
