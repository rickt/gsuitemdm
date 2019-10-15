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
* Lots more that will be highlighted soon

## Status
NOT YET READY FOR PUBLIC USE

## Updates
* 20191015: Basic features all ported into package & working (Admin SDK API, Datastore, Sheets, Data conversion & merging)
* 20191014: Started conversion to go package

## TODO
* Port all MDM device action operations (Approve/Block/Wipe Account/Wipe Device) into package
* Port all search operations into package

