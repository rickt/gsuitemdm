package mdmtool

//
// Types for mdmtool
//

import (
	"github.com/rickt/gsuitemdm"
)

type MDMTool struct {
	GSMDMService gsuitemdm.GSuiteMDMService // G Suite MDM Service
	URLs         MDMToolURLs                // URLs

}

type MDMToolURLs struct {
	ApproveDeviceURL   string `json:"approveurl"`         // URL of Approve cloud function
	BlockDeviceURL     string `json:"blockurl"`           // URL of Block cloud function
	DeleteDeviceURL    string `json:"deleteurl"`          // URL of Delete cloud function
	DirectoryURL       string `json:"directoryurl"`       // URL of Directory cloud function
	SearchDatastoreURL string `json:"searchurl"`          // URL of Search cloud function
	UpdateDatastoreURL string `json:"updatedatastoreurl"` // URL of Update Datastore cloud function
	UpdateSheetURL     string `json:"updatesheeturl"`     // URL of Update Sheet cloud function
	WipeDeviceURL      string `json:"wipeurl"`            // URL of Wipe cloud function
}

// EOF
