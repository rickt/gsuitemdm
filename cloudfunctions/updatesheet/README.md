# gsuitemdm Cloud Function `updatesheet` #

A [cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that updates a Google Sheet with the most recent mobile device data from [Google Datastore](https://cloud.google.com/datastore/).

## HOW-TO Configure `updatesheet` ##
`updatesheet` uses a `.yaml` file containing several environment variables the cloud function reads during app startup. These environment variables point the app to the shared master cloud function configuration and API key that are stored as [Secret Manager secrets](https://cloud.google.com/secret-manager/docs/managing-secrets). An example `.yaml` file for `updatesheet`:

```yaml
APPNAME: updatesheet
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```

## HOW-TO Deploy `updatesheet` ##
```
$ gcloud functions deploy UpdateSheet \
  --runtime go111 \
  --trigger-http \
  --env-vars-file env_updatesheet.yaml 
```

## HOW-TO Use `updatesheet` ##

### API ###
Example expected JSON to update the Google Sheet with the most recent mobile device data from Google Datastore:

```json
{
	"key": "0123456789"
}
```

Example command line using `curl` and the above JSON to update the Google Sheet with the most recent mobile device data from Google Datastore:

```
$ curl -X POST -d '{"key": "0123456789"}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/UpdateSheet
```

It is recommended that Google [Cloud Scheduler](https://cloud.google.com/scheduler/) be used to schedule automatic calls to `updatesheet`, so that your ops' team can have an automatically-updated Google Sheet containing all users' mobile device infomation. For example:

```
$ gcloud scheduler jobs list
ID                   LOCATION     SCHEDULE (TZ)                      TARGET_TYPE  STATE
UpdateSheet          us-central1  */5 * * * * (America/Los_Angeles)  HTTP         ENABLED
```

### mdmtool ###
The Google Sheet can be updated by running `mdmtool`'s `updatesheet` command:

```
$ mdmtool updatesheet
Updating Google Sheet...  done.
updatesheet Success
```

