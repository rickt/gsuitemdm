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

## Data Storage ##
* Mobile device & user data stored in [Google Datastore](https://cloud.google.com/datastore/docs/)
* Configuration/keys stored in Google [Secret Manager](https://cloud.google.com/secret-manager/docs/)
* (Optional) Google Sheet for mobile device tracking

## Configuration ##
All configuration data, API keys and service account domain credentials are stored as secrets in Google [Secret Manager](https://cloud.google.com/secret-manager/docs/). Learn more about [`gsuitemdm` configuration](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions#configuration) or [`gsuitemdm` secrets](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions#configuration-secrets).

## Status
* In production
* Ready for public use
* Docs: 90%

## Pre-Requisites ##
* 1+ G Suite domain(s) using G Suite MDM to manage iOS/Android mobile devices
* GCP project with billing setup

## Brief Setup Notes
A full and complete installation guide to follow, but for now, some brief setup notes & requirements: 

* Setup GCP project 
  * Enable the required APIs (`admin`, `cloudfunctions`, `cloudscheduler`, `datastore`, `logging`, `secretmanager`, `sheets`)
* Setup service account(s) + JSON credentials `foreach` G Suite domain including [G Suite domain-wide delegation](https://developers.google.com/admin-sdk/directory/v1/guides/delegation)
* Create secrets in Secret Manager for: 
  * G Suite domain JSON credentials
  * API key
  * Slack security token
* Grant appropriate scopes to service accounts in the Admin Console
* Setup Google Datastore
* Setup Google Sheet template to track mobile devices

### Google Sheet Setup
1. Make a copy of [this Google Sheet](https://update.url) and save it in Google Drive. Now get the ID of your sheet; this is the part after `https://docs.google.com/spreadsheets/d/` in the sheet's URL but before `/edit`. Add that sheet ID to the main JSON configuration file, `"sheetid": "yourgooglesheetidgoeshere"`
2. Add the email address of the G Suite user who you wish to update the Google sheet as, to the main JSON configuration file, `"sheetwho": "username@yourgsuitedomain.com"`

