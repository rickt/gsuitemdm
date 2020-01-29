# gsuitemdm Cloud Function `blockdevice` #

A [cloud Function](https://cloud.google.com/functions/) component of the [`gsuitemdm`](https://github.com/rickt/gsuitemdm) package that blocks a mobile device using the [Admin SDK](https://developers.google.com/admin-sdk).

## HOW-TO Deploy `blockdevice` ##
```
$ gcloud functions deploy BlockDevice --runtime go111 --trigger-http \
  --env-vars-file env_blockdevice.yaml
```

## HOW-TO Use `blockdevice` ##

### API ###
Example expected JSON to block a device in the domain `foo.com` with IMEI `1234567890987654321`:
```json
{
	"key": "0123456789",
	"action": "block",
	"imei": "1234567890987654321",
	"domain": "foo.com",
	"confirm": true
}
```

Note that if `"confirm": true` is not specified, the device will not be blocked. 

Example command line using `curl` and the above JSON to block the abovementioned device:

```
$ curl -X POST -d \
  '{"key": "0123456789", "action": "block", "imei": "1234567890987654321", "domain": "foo.com", "confirm": true}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/BlockDevice
```

### `mdmtool` ##
Example command line using `mdmtool` to block a device in the domain `foo.com` with IMEI `1234567890987654321`:
```
$ mdmtool block -i 1234567890987654321 -d foo.com
WARNING: Are you sure you want to APPROVE device IMEI=1234567890987654321 in domain foo.com? [y/n]: 
```
