# gsuitemdm Cloud Function `deletedevice` #

A [Cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that deletes a mobile device using the [Admin SDK](https://developers.google.com/admin-sdk).

## HOW-TO Deploy ##
`$ gcloud functions deploy DeleteDevice --runtime go111 --trigger-http --env-vars-file env_deletedevice.yaml`

## API Examples ##
Example test command line that deletes a device in the domain foo.com with IMEI 1234567890987654321:

```
$ curl -X POST -d \
  '{"key": "0123456789", "action": "delete", "imei": "1234567890987654321", "domain": "foo.com", "confirm": true}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/DeleteDevice
```
