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
* `Block`
	* `$ mdmtool block -s SN -d DOMAIN`
* `Delete`
	* `$ mdmtool delete -i IMEI -d DOMAIN`
* `Wipe`
	* `$ mdmtool wipe -s SN -d DOMAIN`


