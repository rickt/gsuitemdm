package gsuitemdm

//
// GSuiteMDM types
//

import (
	"context"
	admin "google.golang.org/api/admin/directory/v1"
)

// G Suite MDM Service main struct type
type GSuiteMDMService struct {
	C             GSuiteMDMConfig         // Main configuration
	Ctx           context.Context         // Context
	DatastoreData []DatastoreMobileDevice // Datastore mobile device data
	SDKData       *admin.MobileDevices    // Admin SDK mobile device data
	SheetData     []DatastoreMobileDevice // Google Sheet mobile device data
}

// G Suite MDM Service config struct type
type GSuiteMDMConfig struct {
	// Required G Suite Admin SDK scope to perform ACTION operations (delete, wipe, block, etc).
	// See SearchScope for more details
	ActionScope string `json:"actionscope"`

	// Default sort order of devices returned by the Admin API query parameter: orderBy.
	// Refer to https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/list
	APIQueryOrderBy string `json:"apiqueryorderby"`

	// Default sort order of devices returned by Cloud Datastore
	DatastoreQueryOrderBy string `json:"datastorequeryorderby"`

	// Global debug mode?
	Debug bool `json:"globaldebug"`

	// Datastore namekey
	DSNamekey string `json:"dsnamekey"`

	// Project ID of the GCP project
	ProjectID string `json:"projectid"`

	// What type of Remote Wipe will we use for the "wipe" command? Possible values are:
	// Refer to https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/action
	RemoteWipeType string `json:"remotewipetype"`

	// Required G Suite Admin API scope to perform SEARCH operations. Since we are using the
	// Mobiledevices: list method of the G Suite Admin API, refer to
	// https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/list
	//
	// Default value for this should be: "https://www.googleapis.com/auth/admin.directory.device.mobile.readonly"
	// and there should likely be no good reason to change it.
	SearchScope string `json:"searchscope"`

	// Default type of search.
	// This refers to the STATUS of the mobile device as seen in the G Suite Admin console.
	// Possible values are:
	//		All
	//		Approved
	//		Pending Approval
	//		Blocked
	//		Account Wiped
	//		Device Wiped
	//		Account Wiping
	//		Device Wiping
	//
	SearchType string `json:"searchtype"`

	// GCP Secret Manager ID of the credentials with necessary permissions to write to the Google Sheet
	SheetCredsID string `json:"sheetcredsid"`

	// ID of the google spreadsheet to update
	SheetID string `json:"sheetid"`

	// Required Sheets API scope to update the Google Sheet
	SheetScope string `json:"sheetscope"`

	// Who to write the spreadsheet as
	SheetWho string `json:"sheetwho"`

	// Time Zone
	TimeZone string `json:"timezone"`

	// Version of gsuitemdm
	Version string `json:"version"`

	// G Suite domains that mdmtool knows about.
	Domains Domains `json:"domains"`
}

// Struct used for configured domains. Just a slice of the domain-specific configuration struct
type Domains []DomainConf

// Specific G Suite domain configuration
type DomainConf struct {
	// Administrative User on this G Suite domain you want the API calls to "run as".
	// This will need to be a user/email address that has Administrator/Super Administrator role
	// in the specific G Suite domain.
	AdminUser string `json:"adminuser"`

	// Immutable Customer ID of G Suite domain
	CustomerID string `json:"companyid"`

	// FQDN of G Suite domain
	DomainName string `json:"domainname"`

	// Credentials Secret for this G Suite domain
	SecretID string `json:"secretid"`
}

// EOF
