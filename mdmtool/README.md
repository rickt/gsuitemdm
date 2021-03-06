# mdmtool

## Overview
A command line utility enabling fast & easy MDM ops on [G Suite MDM-protected mobile devices](https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices).

## Requirements
Requires all [G Suite MDM Cloud Functions](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions) to be setup, properly deployed and working. Change [mdmtool configuration file](https://github.com/rickt/gsuitemdm/blob/master/mdmtool/mdmtool_conf_example.json) to point to your own Cloud Function URLs and your own API key. 

## Actions
Actions perform the requested [Admin SDK administrative operation](https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/action) on a mobile device. 

All actions require `-i IMEI` or `-s SN` as well as `-d DOMAIN`, e.g.
* `Approve` 
	* `$ mdmtool approve -i IMEI -d DOMAIN`
* `Block` 
	* `$ mdmtool block -s SN -d DOMAIN`
* `Delete` 
	* `$ mdmtool delete -i IMEI -d DOMAIN`
* `Wipe` 
	* `$ mdmtool wipe -s SN -d DOMAIN`

All actions require a valid (Y/N) confirmation response before being executed, e.g.
```
$ mdmtool wipe -i 123456789098765 -d foo.com
WARNING: Are you sure you want to WIPE device IMEI=123456789098765 in domain foo.com? [y/n]: 
```

### Action Types
| Action  | What it does                 | Details on what it does                                                              |
|---------|------------------------------|--------------------------------------------------------------------------------------|
| `Approve` | Approves a mobile device     | Allows a user to sign into G Suite on their mobile device                            |
| `Block`   | Blocks a mobile device       | Remotely log out signed-in users, disable ability to login to mobile device          |
| `Delete`  | Deletes a mobile device      | Removes a device from MDM; use only when replacing a mobile device with a new one    |
| `Wipe`    | Remote-wipes a mobile device | Forcibly remove all data & content from a device; device returns to factory settings |

See the [Mobiledevices: action Admin SDK docs](https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/action) for full details on G Suite MDM administrative actions. 

## Directory
Search for user phone numbers.
```
$ mdmtool dir -n doe
----------------------+----------------+------------------------------------------
Name                  | Phone Number   | Email 
----------------------+----------------+------------------------------------------
Jane Doe              | (213) 555-1212 | jane@foo.com
John Doe              | (323) 555-1212 | john@bar.com
----------------------+----------------+------------------------------------------
Search returned 2 results.
```
* Search for phone number using device owner name:
	* `$ mdmtool dir -n john`
* Search for phone number using device owner email address:
	* `$ mdmtool dir -e john@bar.com`

## Search
Search for mobile device data.

```
$ mdmtool search -n doe
----------------------+------------------+----------------+------------------+-----------------+---------------+--------------------+---------------
Domain                | Model            | Phone Number   | Serial #         | IMEI            | Status        | Last Sync          | Owner
----------------------+------------------+----------------+------------------+-----------------+---------------+--------------------+---------------
foo.com               | iPhone 5S        | (213) 555-1212 | ABC123ABC123     | 123456789098765 | APPROVED      | 1 hour ago         | Jane Doe
bar.com               | iPhone XR        | (323) 555-1212 | ZX81TRS80C64     | 234567890987654 | BLOCKED       | 11 hours ago       | John Doe
----------------------+------------------+----------------+------------------+-----------------+---------------+--------------------+---------------
Search returned 2 results.
```
* Search using device owner name:
	* `$ mdmtool search -n john`
* Search using device owner email address:
	* `$ mdmtool search -e foo@bar.com`
* Search using device IMEI:
	* `$ mdmtool search -i 123456789098765`
* Search using device serial number:
	* `$ mdmtool search -n ZX81TRS80C64`
* Search using notes (stored in a device tracking Google Sheet):
	* `$ mdmtool search -o lost`
* Search using phone number (stored in a device tracking Google Sheet):
	* `$ mdmtool search -p 2135551212`
* Search using device status:
	* `$ mdmtool search -t BLOCKED`

## Updates
* `Update Datastore`
	* `$ mdmtool udpatedb`
* `Update Google Sheet`
	* `$ mdmtool updateshet`

### Update Types
| Update Type       | Details on what it does                                                                      |
|-------------------|----------------------------------------------------------------------------------------------|
| `updatedatastore` | Updates Datastore with fresh data from Admin SDK for all devices, merge w/Google Sheet data, save to Datastore |
| `updatesheet`    | Updates Google Sheet with fresh data from Datastore                                          |


