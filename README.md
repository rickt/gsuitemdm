# gsuitemdm
`gsuitemdm` is a Go package that eases the management of iOS or Android mobile devices in multiple G Suite domains that use [G Suite MDM](https://support.google.com/a/answer/1734200?hl=en) to secure their mobile devices.

`gsuitemdm` provides:
* Multiple, easy to use, secure mobile device management APIs deployed as [cloud functions](https://cloud.google.com/functions/docs/) to help you quickly manage many mobile devices 
* A command line tool ([`mdmtool`](https://github.com/rickt/gsuitemdm/tree/master/mdmtool) allowing for easy command line mobile device management
* Persistent data storage in Google [Datastore](https://cloud.google.com/datastore/)
* Configuration & credentials stored securely in Google [Secret Manager](https://cloud.google.com/secret-manager/docs/)

## Additional Features ##
* Additional APIs such as a [phone directory](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/directory) and [Slack `/phone` command](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/directory)
* Uses [GCP service accounts](https://developers.google.com/identity/protocols/OAuth2ServiceAccount) and G Suite [domain-wide delegation authority](https://gsuite-developers.googleblog.com/2012/11/domain-wide-delegation-of-authority-and.html)
* Supports multiple G Suite domains, easy & shared configuration across all components
* Quickly and easily Approve/Block/Delete/Wipe/Search for MDM-protected devices across multiple G Suite domains
* Generate an auto-updating [Google Sheet](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/updatesheet) so your ops team can track all mobile devices across multiple G Suite domains

## Use-Cases ##

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

