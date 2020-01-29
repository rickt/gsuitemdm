# gsuitemdm
GSuiteMDM is a Go package that eases the management of iOS or Android mobile devices in G Suite domains that use G Suite MDM to secure their mobile devices.

GSuiteMDM provides:
* Multiple, easy to use, secure APIs to help you quickly manage many mobile devices 
	* G Suite MDM is a high-layer wrapper around the [G Suite Admin SDK API](https://developers.google.com/admin-sdk)
	* Supports multiple G Suite domains, E-Z configuration
	* Uses [GCP service accounts](https://developers.google.com/identity/protocols/OAuth2ServiceAccount) and G Suite [domain-wide delegation authority](https://gsuite-developers.googleblog.com/2012/11/domain-wide-delegation-of-authority-and.html)
	* Mobile device data stored in Google Datastore
	* Quickly and easily Approve/Block/Delete/List/Search for/Wipe MDM-protected devices in multiple domains
		* Build your own web app by calling the GSuiteMDM APIs, or
		* Use the built-in command line tool `mdmtool` (coming soon)
	* Search for mobile devices based on:
		* Owner G Suite account full name or email address
		* Mobile device current G Suite MDM Status (APPROVED, PENDING, BLOCKED, WIPING, etc)
		* Mobile device phone number (TODO: expand notes on this + Google Sheet)
		* Mobile device IMEI or Serial Number
		* Any notes associated with the device (TODO: expand notes on this + Google Sheet)
		* G Suite domain to which the device belongs
	* API endpoints are simple, easy to deploy and manage, lightweight GCP [Cloud Functions](https://cloud.google.com/functions/). Authentication is done by key, e.g. 
```
$ curl -X POST -d '{"key": "$KEY", "qtype": "name", "q": "john smith"}' https://$APIURL/SearchApiEndpoint
[
   {
      "Color": "black",
      "CompromisedStatus": "No compromise detected",
      "Domain": "foo.com",
      "DeveloperMode": false,
      "Email": "jsmith@foo.com",
      "IMEI": "01234567890987654321",
      "Model": "iPhone XR",
      "Name": "John Smith",
	[ snipped for brevity ]
      "USBADB": false,
      "WifiMac": "aa:bb:cc:dd:ee:ff"
   }
]
```

## Features
Out of the box, GSuiteMDM offers 3 main things to help G Suite MDM administrators:
1. A simple to use (but secure) mobile device management API that lets you manage devices in multiple domains without pulling your hair out dealing with the G Suite Admin SDK. Cloud Functions/APIs available:
	* ApproveDevice
	* BlockDevice
	* DeleteDevice
	* Directory 
	* SearchDatastore
	* UpdateDatastore
	* UpdateSheet
	* WipeDevice
2. Tracking Google Sheet
	* (TODO)
2. tools (command line tool `mdmtool` and auto-updated mobile device tracking Google Sheet to manage G Suite MDM-protected mobile devices without having to use the G Suite [Admin Console](https://admin.google.com/). 

## Configuration ##
All configuration data, API keys and JSON G Suite domain credentials are stored as secrets in Google [Secret Manager](https://cloud.google.com/secret-manager/docs/). Learn more about [`gsuitemdm` configuration](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions#configuration) or [`gsuitemdm` secrets](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions#configuration-secrets).

## Status
READY FOR PUBLIC USE

## Updates
* 20191022: All features tested, working, ready for public use
* 20191021: Converted to POST for security reasons
* 20191015: Basic features all ported into package & working (Admin SDK API, Datastore, Sheets, Data conversion & merging)
* 20191014: Started conversion to go package

## TODO
* [DONE] Add PhoneNumber API cloud function
* [DONE] Port all MDM device action operations (Approve/Block/Wipe Account/Wipe Device) into package
* [DONE] Port all search operations into package

## Setup Notes

### GCP Project Setup
gsuitemdm needs a GCP project to run inside. Pick an existing GCP project, or create a new one. Either way, add it to the main JSON configuration file, `"projectid": "yourgcpprojectname"`. 

### Per-G Suite domain Credentials Setup
Docs coming. 

### Datastore Setup
Docs coming.

### Google Sheet Setup
1. Make a copy of [this Google Sheet](https://update.url) and save it in Google Drive. Now get the ID of your sheet; this is the part after `https://docs.google.com/spreadsheets/d/` in the sheet's URL but before `/edit`. Add that sheet ID to the main JSON configuration file, `"sheetid": "yourgooglesheetidgoeshere"`
2. Add the email address of the G Suite user who you wish to update the Google sheet as, to the main JSON configuration file, `"sheetwho": "username@yourgsuitedomain.com"`

