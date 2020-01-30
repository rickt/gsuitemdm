package main

//
// Types for mdmtool
//

import ()

type MDMTool struct {
	Config MDMToolConfig // MDMTool configuration

}

type MDMToolConfig struct {
	APIKey             string `json:"apikey"`             // G Suite MDM API Key
	ApproveDeviceURL   string `json:"approvedeviceurl"`   // URL of Approve Device cloud function
	BlockDeviceURL     string `json:"blockdeviceurl"`     // URL of Block Device cloud function
	DeleteDeviceURL    string `json:"deletedeviceurl"`    // URL of Delete Device cloud function
	DirectoryURL       string `json:"directoryurl"`       // URL of Directory cloud function
	SearchDatastoreURL string `json:"searchdatastoreurl"` // URL of Search Device cloud function
	UpdateDatastoreURL string `json:"updatedatastoreurl"` // URL of Update Datastore cloud function
	UpdateSheetURL     string `json:"updatesheeturl"`     // URL of Update Sheet cloud function
	WipeDeviceURL      string `json:"wipedeviceurl"`      // URL of Wipe Device cloud function
}

// GSuiteMDM URLs
type MDMToolURLs struct {
	ApproveDeviceURL   string `json:"approvedeviceurl"`   // URL of Approve Device cloud function
	BlockDeviceURL     string `json:"blockdeviceurl"`     // URL of Block Device cloud function
	DeleteDeviceURL    string `json:"deletedeviceurl"`    // URL of Delete Device cloud function
	DirectoryURL       string `json:"directoryurl"`       // URL of Directory cloud function
	SearchDatastoreURL string `json:"searchdatastoreurl"` // URL of Search Device cloud function
	UpdateDatastoreURL string `json:"updatedatastoreurl"` // URL of Update Datastore cloud function
	UpdateSheetURL     string `json:"updatesheeturl"`     // URL of Update Sheet cloud function
	WipeDeviceURL      string `json:"wipedeviceurl"`      // URL of Wipe Device cloud function
}

// EOF
