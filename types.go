package gsuitemdm

//
// GSuiteMDM types
//

import (
	"cloud.google.com/go/logging"
	"context"
)

// G Suite MDM Service main struct type
type GSuiteMDMService struct {
	C             GSuiteMDMConfig         // Main configuration
	Ctx           context.Context         // Context
	Log           *logging.Client         // Stackdriver (GCP) log
	SheetData     []DatastoreMobileDevice // Google Sheet mobile device data
	DatastoreData []DatastoreMobileDevice // Datastore mobile device data
}

// G Suite MDM Service config struct type
type GSuiteMDMConfig struct {
	// Required G Suite Admin SDK scope to perform ACTION operations (delete, wipe, block, etc).
	// See SearchScope for more details
	ActionScope string `json:"actionscope"`

	// Global debug mode?
	Debug bool `json:"globaldebug"`

	// Datastore namekey
	DSNamekey string `json:"dsnamekey"`

	// Required G Suite Admin API scope to perform SEARCH operations. Since we are using the
	// Mobiledevices: list method of the G Suite Admin API, refer to
	// https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/list
	//
	// Default value for this should be: "https://www.googleapis.com/auth/admin.directory.device.mobile.readonly"
	// and there should likely be no good reason to change it.
	SearchScope string `json:"searchscope"`

	// Required Sheets API scope to update the Google Sheet
	SheetScope string `json:"sheetscope"`

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

	// G Suite domains that mdmtool knows about.
	Domains Domains `json:"domains"`

	// Default sort order of devices returned by the Admin API query parameter: orderBy.
	// Refer to https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/list
	APIQueryOrderBy string `json:"apiqueryorderby"`

	// Default sort order of devices returned by Cloud Datastore
	DatastoreQueryOrderBy string `json:"datastorequeryorderby"`

	// Version of gsuitemdm
	Version string `json:"version"`

	// JSON credentials file for user with necessary permissions to write the Google Sheet
	SheetCreds string `json:"sheetcreds"`

	// ID of the google spreadsheet to update
	SheetID string `json:"sheetid"`

	// Who to write the spreadsheet as
	SheetWho string `json:"sheetwho"`

	// What type of Remote Wipe will we use for the "wipe" command? Possible values are:
	// Refer to https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/action
	RemoteWipeType string `json:"remotewipetype"`

	// Project ID of the GCP project
	ProjectID string `json:"projectid"`
}

// Struct used for configured domains. Just an array of the domain-specific configuration struct
type Domains []DomainConf

// Specific G Suite domain configuration
type DomainConf struct {
	// FQDN of G Suite domain
	DomainName string `json:"domainname"`

	// Immutable Customer ID of G Suite domain
	CustomerID string `json:"companyid"`

	// Administrative User on this G Suite domain you want the API calls to "run as".
	// This will need to be a user/email address that has Administrator/Super Administrator role
	// in the specific G Suite domain.
	AdminUser string `json:"adminuser"`

	// JSON credentials file for this domain's GCP service account.
	// This GCP service account must be:
	// 1) Granted Domain-Wide Delegation authority. See https://developers.google.com/admin-sdk/directory/v1/guides/delegation
	// 2) The ClientId of this service account must be granted access to the following scope:
	//    https://www.googleapis.com/auth/admin.directory.device.mobile.readonly
	//    Grant scope access in G Suite Admin Console Admin Console –> Security –> Advanced Settings –>
	//			--> Authentication –> Manage API Client Access
	CredentialsFile string `json:"credentialsfile"`
}

// Example configuration file:

//     {
//       "scope": "https://www.googleapis.com/auth/admin.directory.device.mobile.readonly",
//       "searchtype": "all",
//       "orderby": "name",
//       "version": "0.3",
//       "domains": [
//       	{
//       		"domainname": "foo.com",
//     	  	"companyid": "A0123ABC",
//     		  "adminuser": "adminuser@foo.com",
//     		  "credentialsfile": "conf/credentials_foo.com.json"
//       	},
//       	{
//       		"domainname": "bar.com",
//       		"companyid": "B0123ABC",
//       		"adminuser": "adminuser@bar.com",
//       		"credentialsfile": "conf/credentials_bar.com.json"
//       	},
//       	{
//       		"domainname": "baz.com",
//       		"companyid": "C0123ABC",
//       		"adminuser": "adminuser@baz.com",
//       		"credentialsfile": "conf/credentials_baz.com.json"
//       	}
//       ]
//     }

// EOF
