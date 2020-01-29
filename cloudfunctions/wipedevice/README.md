# gsuitemdm Cloud Function `wipedevice` #

A [cloud Function](https://cloud.google.com/functions/) component of the [`gsuitemdm`](https://github.com/rickt/gsuitemdm) package that [wipes a mobile device](https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/delete) using the [Admin SDK](https://developers.google.com/admin-sdk).

## HOW-TO Configure `wipedevice` ##
`wipedevice` uses a `.yaml` file containing several environment variables the cloud function reads during app startup. These environment variables point the app to the shared master cloud function configuration and API key that are stored as [Secret Manager secrets](https://cloud.google.com/secret-manager/docs/managing-secrets). An example `.yaml` file for `wipedevice`:

```yaml
APPNAME: wipedevice
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```

## HOW-TO Deploy `wipedevice` ##
```
$ gcloud functions deploy WipeDevice --runtime go111 --trigger-http \
  --env-vars-file env_wipedevice.yaml
```

## HOW-TO Use `wipedevice` ##

### API ###
Example expected JSON to wipe a device in the domain `foo.com` with IMEI `1234567890987654321`:
```json
{
	"action": "wipe",
	"confirm": true,
	"domain": "foo.com",
	"imei": "1234567890987654321",
	"key": "0123456789"
}
```

Note that if `"confirm": true` is not specified, the device will not be wiped. 

Example command line using `curl` and the above JSON to wipe the abovementioned device:

```
$ curl -X POST -d \
  '{"key": "0123456789", "action": "wipe", "imei": "1234567890987654321", "domain": "foo.com", "confirm": true}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/WipeDevice
```

### `mdmtool` ##
Example command line using `mdmtool` to wipe a device in the domain `foo.com` with IMEI `1234567890987654321`:
```
$ mdmtool wipe -i 1234567890987654321 -d foo.com
WARNING: Are you sure you want to WIPE device IMEI=1234567890987654321 in domain foo.com? [y/n]: 
```
