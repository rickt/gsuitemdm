# gsuitemdm
A go package to ease the management and operations of G Suite MDM-protected corporate mobile devices. 

If you:

* manage 1+ G Suite domains
* and manage your corporate mobile devices using [G Suite MDM](https://support.google.com/a/answer/1734200?hl=en)

then gsuitemdm might be right up your alley. 

gsuitemdm offers a framework and tools to manage G Suite MDM-protected mobile devices without having to use the G Suite [Admin Console](https://admin.google.com/). 

## Features
* Stores mobile device state in Google Datastore for speed (the G Suite [Admin SDK](https://developers.google.com/admin-sdk) is quite slow) and resiliency
* Easily create an automatically-updating Google Sheet to track all your G Suite domains' mobile devices
* API Endpoints deployed as (Cloud Functions)[https://cloud.google.com/functions/] offering functionality such as:
	* Update Datastore with fresh mobile device data from the Admin SDK
	* Update Google Sheet with mobile device data from Google Datastore
	* Search Datastore for mobile device(s) based on criteria such as owner name, email address, IMEI/SN, etc

## Status
ALMOST READY FOR PUBLIC USE

## Updates
* 20191021: Converted to POST for security reasons
* 20191015: Basic features all ported into package & working (Admin SDK API, Datastore, Sheets, Data conversion & merging)
* 20191014: Started conversion to go package

## TODO
* Port all MDM device action operations (Approve/Block/Wipe Account/Wipe Device) into package
* Port all search operations into package

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

