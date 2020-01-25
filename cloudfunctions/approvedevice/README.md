# gsuitemdm Cloud Function `approvedevice` #

A [Cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that approves a mobile device using the [Admin SDK](https://developers.google.com/admin-sdk).

## HOW-TO Deploy ##
`$ gcloud functions deploy ApproveDevice --runtime go111 --trigger-http --env-vars-file env_approvedevice.yaml`

## API Examples ##
Example test command line that approves a device in the domain foo.com with IMEI 1234567890987654321:

```
$ curl -X POST -d \
  '{"key": "0123456789", "action": "approve", "imei": "1234567890987654321", "domain": "foo.com", "confirm": true}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/ApproveDevice
```

