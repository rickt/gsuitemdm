# gsuitemdm Cloud Function `deletedevice` #

A [cloud Function](https://cloud.google.com/functions/) component of the [`gsuitemdm`](https://github.com/rickt/gsuitemdm) package that [deletes a mobile device](https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/delete) using the [Admin SDK](https://developers.google.com/admin-sdk).

## HOW-TO Deploy `deletedevice` ##
```
$ gcloud functions deploy DeleteDevice --runtime go111 --trigger-http \
  --env-vars-file env_deletedevice.yaml
```

## HOW-TO Use `deletedevice` ##

### API ###
Example expected JSON to delete a device in the domain `foo.com` with IMEI `1234567890987654321`:
```json
{
	"key": "0123456789",
	"action": "delete",
	"imei": "1234567890987654321",
	"domain": "foo.com",
	"confirm": true
}
```

Note that if `"confirm": true` is not specified, the device will not be deleted. 

Example command line using `curl` and the above JSON to delete the abovementioned device:

```
$ curl -X POST -d \
  '{"key": "0123456789", "action": "delete", "imei": "1234567890987654321", "domain": "foo.com", "confirm": true}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/ApproveDevice
```

### `mdmtool` ##
Example command line using `mdmtool` to delete a device in the domain `foo.com` with IMEI `1234567890987654321`:
```
$ mdmtool delete -i 1234567890987654321 -d foo.com
WARNING: Are you sure you want to DELETE device IMEI=1234567890987654321 in domain foo.com? [y/n]: 
```
