# gsuitemdm Cloud Function `blockdevice` #

A [cloud Function](https://cloud.google.com/functions/) component of the [`gsuitemdm`](https://github.com/rickt/gsuitemdm) package that [blocks a mobile device](https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/action) using the [Admin SDK](https://developers.google.com/admin-sdk).

## HOW-TO Configure `blockdevice` ##
`blockdevice` uses a `.yaml` file containing several environment variables the cloud function reads during app startup. These environment variables point the app to the shared master cloud function configuration and API key that are stored as [Secret Manager secrets](https://cloud.google.com/secret-manager/docs/managing-secrets). An example `.yaml` file for `blockdevice`:

```yaml
APPNAME: blockdevice
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```

## HOW-TO Deploy `blockdevice` ##
```
$ gcloud functions deploy BlockDevice \
  --runtime go111 \
  --trigger-http \
  --env-vars-file env_blockdevice.yaml
```

## HOW-TO Use `blockdevice` ##

### API ###
Example expected JSON to block a device in the domain `foo.com` with IMEI `111111111111111`:
```json
{
	"action": "block",
	"confirm": true,
	"domain": "foo.com",
	"imei": "111111111111111",
	"key": "0123456789"
}
```

Note that if `"confirm": true` is not specified, the device will not be blocked. 

Example command line using `curl` and the above JSON to block the abovementioned device:

```
$ curl -X POST -d \
  '{"key": "0123456789", "action": "block", "imei": "111111111111111", "domain": "foo.com", "confirm": true}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/BlockDevice
```

### `mdmtool` ##
Example command line using `mdmtool` to block a device in the domain `foo.com` with IMEI `111111111111111`:
```
$ mdmtool block -i 111111111111111 -d foo.com
WARNING: Are you sure you want to BLOCK device IMEI=111111111111111 in domain foo.com? [y/n]: 
```
