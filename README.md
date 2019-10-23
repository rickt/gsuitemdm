# gsuitemdm
GSuiteMDM is a higher-layer API framework to ease the management of iOS or Android mobile devices in G Suite domains that use G Suite MDM to secure their mobile devices. 

GSuiteMDM provides:
* Multiple, easy to use, secure APIs to help you quickly manage many mobile devices 
	* A high-layer wrapper around the [G Suite Admin SDK API](https://developers.google.com/admin-sdk)
	* Persistent data storage in cheap/fast/resilient Google Datastore
	* Quickly and easily Approve/Block/Delete/List/Search for/Wipe MDM-protected devices in multiple domains
	* Search for or list all devices based on:
		* Owner name or email address
		* Device current G Suite MDM Status (APPROVED, PENDING, BLOCKED, WIPING, etc)
		* Device phone number
		* Device IMEI or Serial Number
		* Any notes associated with the device
		* G Suite domain to which the device belongs
	* API endpoints are simple, easy to deploy and manage, lightweight GCP [Cloud Functions](https://cloud.google.com/functions/)

## Features
Out of the box, GSuiteMDM offers 3 main things to help G Suite MDM administrators:
1. A simple to use (but secure) mobile device management API that lets you manage devices in multiple domains without pulling your hair out dealing with the G Suite Admin SDK. Cloud Functions/APIs available:
	* ApproveDevice
	* BlockDevice
	* DeleteDevice
	* Directory / Phone Number Search
	* Search (for a mobile device)
	* Update Datastore
	* Automatically update a mobile device tracking Google Sheet
	* WipeDevice
2. Tracking Google Sheet
2. tools (command line tool `mdmtool` and auto-updated mobile device tracking Google Sheet to manage G Suite MDM-protected mobile devices without having to use the G Suite [Admin Console](https://admin.google.com/). 

## Status
READY FOR PUBLIC USE

## Updates
* 20191022: All features tested, working, ready for public use
* 20191021: Converted to POST for security reasons
* 20191015: Basic features all ported into package & working (Admin SDK API, Datastore, Sheets, Data conversion & merging)
* 20191014: Started conversion to go package

## TODO
* Add PhoneNumber API cloud function
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

