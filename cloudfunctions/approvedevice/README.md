# gsuitemdm Cloud Function `approvedevice` #

A [cloud Function](https://cloud.google.com/functions/) component of the [`gsuitemdm`](https://github.com/rickt/gsuitemdm) package that [approves a mobile device](https://developers.google.com/admin-sdk/directory/v1/reference/mobiledevices/action) using the [Admin SDK](https://developers.google.com/admin-sdk).

## HOW-TO Configure `approvedevice` ##
`approvedevice` uses a `.yaml` file containing several environment variables the cloud function reads during app startup. These environment variables point the app to the shared master cloud function configuration and API key that are stored as [Secret Manager secrets](https://cloud.google.com/secret-manager/docs/managing-secrets). An example `.yaml` file for `approvedevice`:

```yaml
APPNAME: approvedevice
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```

## HOW-TO Deploy `approvedevice` ##
```
$ gcloud functions deploy ApproveDevice \
  --runtime go111 \
  --trigger-http \
  --env-vars-file env_approvedevice.yaml
```

## HOW-TO Use `approvedevice` ##

### API ###
Example expected JSON to approve a device in the domain `foo.com` with IMEI `111111111111111`:
```json
{
	"action": "approve",
	"confirm": true,
	"domain": "foo.com",
	"imei": "111111111111111",
	"key": "0123456789"
}
```

Note that if `"confirm": true` is not specified, the device will not be approved. 

Example command line using `curl` and the above JSON to approve the abovementioned device:

```
$ curl -X POST -d \
  '{"key": "0123456789", "action": "approve", "imei": "111111111111111", "domain": "foo.com", "confirm": true}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/ApproveDevice
```

### `mdmtool` ##
Example command line using `mdmtool` to approve a device in the domain `foo.com` with IMEI `111111111111111`:
```
$ mdmtool approve -i 111111111111111 -d foo.com
WARNING: Are you sure you want to APPROVE device IMEI=111111111111111 in domain foo.com? [y/n]: 
```
