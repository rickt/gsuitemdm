# gsuitemdm Cloud Function `approvedevice` #

A [Cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that approves a mobile device using the [Admin SDK](https://developers.google.com/admin-sdk).

## HOW-TO Deploy ##
```
$ gcloud functions deploy ApproveDevice --runtime go111 --trigger-http \
  --env-vars-file env_approvedevice.yaml
```

## HOW-TO Use `approvedevice` ##

### API ###
Example expected JSON to approve a device in the domain `foo.com` with IMEI `1234567890987654321`:
```json
{
	"key": "0123456789",
	"action": "approve",
	"imei": "1234567890987654321",
	"domain": "foo.com",
	"confirm": true
}
```

Note that if `"confirm": true` is not specified, the device will not be approved. 

Example command line using `curl` and the above JSON to approve the abovementioned device:

```
$ curl -X POST -d \
  '{"key": "0123456789", "action": "approve", "imei": "1234567890987654321", "domain": "foo.com", "confirm": true}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/ApproveDevice
```

### `mdmtool` ##
Example command line using `mdmtool` to approve a device in the domain `foo.com` with IMEI `1234567890987654321`:
```
$ mdmtool approve -i 1234567890987654321 -d foo.com
WARNING: Are you sure you want to APPROVE device IMEI=1234567890987654321 in domain foo.com? [y/n]: 
```
