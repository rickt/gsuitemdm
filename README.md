# gsuitemdm
`gsuitemdm` is a Go package that eases the management of iOS or Android mobile devices in G Suite domains that use [G Suite MDM](https://support.google.com/a/answer/1734200?hl=en) to secure their mobile devices.

`gsuitemdm` provides:
* Multiple, easy to use, secure mobile device management APIs deployed as [cloud functions](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/) to help you quickly manage many mobile devices 
* A command line tool ([`mdmtool`](https://github.com/rickt/gsuitemdm/tree/master/mdmtool)) allowing for easy command line mobile device management
* Mobile device & user data stored in [Google Datastore](https://cloud.google.com/datastore/docs/)
* Configuration, keys & credentials stored securely as secrets in Google [Secret Manager](https://cloud.google.com/secret-manager/docs/)

Basically, `gsuitemdm` gives you:
* A much more convenient API interface to the [G Suite Admin SDK](https://developers.google.com/admin-sdk)
* Ability to script MDM operations or use a [CLI tool](https://github.com/rickt/gsuitemdm/tree/master/mdmtool) instead of the Admin Console
* Handy 'quality of life' bits & pieces such as a [phone directory](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/directory) API and [Slack `/phone` command](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/directory)

## Additional Features ##
* Securely uses [GCP service accounts](https://developers.google.com/identity/protocols/OAuth2ServiceAccount), GCP [IAM roles](https://cloud.google.com/iam/docs/overview) and G Suite [domain-wide delegation authority](https://gsuite-developers.googleblog.com/2012/11/domain-wide-delegation-of-authority-and.html)
* Supports multiple G Suite domains with easy (and shared!) configuration across all components
  * G Suite domains do not need to be under the [same G Suite account](https://support.google.com/a/answer/182081?hl=en)
* Quickly and easily perform actions (Approve/Block/Delete/Wipe/Search for) on MDM-protected devices across multiple G Suite domains
* Generate an auto-updating [Google Sheet](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/updatesheet) so your ops team can track all mobile devices across multiple G Suite domains
* Structured application logs in [Stackdriver](https://cloud.google.com/logging/)

## Use-Cases ##
* G Suite administrators managing multiple mobile devices in multiple G Suite domains spread across multiple G Suite organizational accounts
* Programmatically perform administrative actions on G Suite MDM-protected mobile devices 
  * Generate an on-call list using the [`directory` API](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/directory)
* Quickly and easily approve/block/wipe mobile devices in the [command line](https://github.com/rickt/gsuitemdm/tree/master/mdmtool) without logging into the G Suite Admin Console

## Status
* In production
* Ready for public use
* Docs: 95%

## Configuration ##
All configuration data, API keys and service account domain credentials are stored as secrets in Google [Secret Manager](https://cloud.google.com/secret-manager/docs/). Learn more about [`gsuitemdm` configuration](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions#configuration) or [`gsuitemdm` secrets](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions#configuration-secrets).

Read all about configuration in the [`gsuitemdm` setup docs](https://github.com/rickt/gsuitemdm/blob/master/docs/SETUP.md).

## Pre-Requisites ##
* 1+ G Suite domain(s) using G Suite MDM to manage iOS/Android mobile devices
* GCP project with billing setup

## Brief Setup Notes
Full setup documenation is [available here](https://github.com/rickt/gsuitemdm/blob/master/docs/SETUP.md).

## TODO ##

