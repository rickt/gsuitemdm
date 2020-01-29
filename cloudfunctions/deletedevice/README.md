# gsuitemdm Cloud Function `deletedevice` #

A [cloud Function](https://cloud.google.com/functions/) component of the [`gsuitemdm`](https://github.com/rickt/gsuitemdm) package that [deletes a mobile device](https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/delete) using the [Admin SDK](https://developers.google.com/admin-sdk).

## HOW-TO Configure `deletedevice` ##
`deletedevice` uses a `.yaml` file containing several environment variables the cloud function reads during app startup. These environment variables point the app to the shared master cloud function configuration and API key that are stored as [Secret Manager secrets](https://cloud.google.com/secret-manager/docs/managing-secrets). An example `.yaml` file for `deletedevice`:

```yaml
APPNAME: deletedevice
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```

## HOW-TO Deploy `deletedevice` ##
```
$ gcloud functions deploy DeleteDevice \
  --runtime go111 \
  --trigger-http \
  --env-vars-file env_deletedevice.yaml
```

## HOW-TO Use `deletedevice` ##

### API ###
Example expected JSON to delete a device in the domain `foo.com` with IMEI `111111111111111`:
```json
{
	"action": "delete",
	"confirm": true,
	"domain": "foo.com",
	"imei": "111111111111111",
	"key": "0123456789"
}
```

Note that if `"confirm": true` is not specified, the device will not be deleted. 

Example command line using `curl` and the above JSON to delete the abovementioned device:

```
$ curl -X POST -d \
  '{"key": "0123456789", "action": "delete", "imei": "111111111111111", "domain": "foo.com", "confirm": true}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/DeleteDevice
```

### `mdmtool` ##
Example command line using `mdmtool` to delete a device in the domain `foo.com` with IMEI `111111111111111`:
```
$ mdmtool delete -i 111111111111111 -d foo.com
WARNING: Are you sure you want to DELETE device IMEI=111111111111111 in domain foo.com? [y/n]: 
```
