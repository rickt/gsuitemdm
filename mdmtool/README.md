# mdmtool
A command line utility enabling fast & easy MDM ops on G Suite MDM-protected mobile devices. 

## Examples

### Search
```
$ mdmtool search -n doe
----------------------+------------------+----------------+------------------+-----------------+---------------+--------------------+---------------
Domain                | Model            | Phone Number   | Serial #         | IMEI            | Status        | Last Sync          | Owner
----------------------+------------------+----------------+------------------+-----------------+---------------+--------------------+---------------
foo.com               | iPhone 5S        | (213) 555-1212 | ABC123ABC123     | 012345678901234 | APPROVED      | 1 hour ago         | Jane Doe
bar.com               | iPhone XR        | (323) 555-1212 | DEF456DEF456     | 357341093918792 | BLOCKED       | 11 hours ago       | John Doe
----------------------+------------------+----------------+------------------+-----------------+---------------+--------------------+---------------
```

### Actions
Actions perform the requested Admin SDK action. All actions require `-i IMEI` or `-s SN` as well as `-d DOMAIN`, e.g.
* `Approve` - Approves a mobile device. Device must be in PENDING or BLOCKED state to be approved. 
	* `$ mdmtool approve -i IMEI -d DOMAIN`
* `Block` - Blocks a mobile device by immediately signing out the G Suite account from all Google apps on the device. Device must be in APPROVED or PENDING state to be blocked. 
	* `$ mdmtool block -s SN -d DOMAIN`
* `Delete` - Deletes a mobile device. This action performs a BLOCK and forces reinstallation of the MDM profile on the device. Device must be in APPROVED, BLOCKED, PENDING, WIPED or WIPING states to be deleted. 
	* `$ mdmtool delete -i IMEI -d DOMAIN`
* `Wipe` - Remote-wipes a mobile device. This action removes ALL DATA AND CONTENTS of the device. Device can be in any state to be wiped. 
	* `$ mdmtool wipe -s SN -d DOMAIN`

A confirmation dialog before any action requires a (Y/N) response. 

| Action  | What it does                 | Allowed States                                     | Example Usage                         |
|---------|------------------------------|----------------------------------------------------|---------------------------------------|
| Approve | Approves a mobile device     | BLOCKED PENDING PROVISIONING                       | `$ mdmtool approve -i IMEI -d DOMAIN` |
| Block   | Blocks a mobile device       | APPROVED PENDING                                   | `$ mdmtool block -i IMEI -d DOMAIN`   |
| Delete  | Deletes a mobile device      | APPROVED BLOCKED PENDING PROVISIONING WIPED WIPING |                                       |
| Wipe    | Remote-wipes a mobile device |                                                    |                                       |

