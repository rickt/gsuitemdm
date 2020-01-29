# gsuitemdm Cloud Function `updatedatastore` #

A [cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that updates mobile device data in [Google Datastore](https://cloud.google.com/datastore/) with the latest mobile device data from the [Admin SDK](https://developers.google.com/admin-sdk) as well as any local device-specific notes that may be stored in the Google Sheet. 

## HOW-TO Configure `updatedatastore` ##
`updatedatastore` uses a `.yaml` file containing several environment variables the cloud function reads during app startup. These environment variables point the app to the shared master cloud function configuration and API key that are stored as [Secret Manager secrets](https://cloud.google.com/secret-manager/docs/managing-secrets). An example `.yaml` file for `updatedatastore`:

```yaml
APPNAME: updatedatastore
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```

## HOW-TO Deploy `updatedatastore` ##
```
$ gcloud functions deploy UpdateDatastore --runtime go111 --trigger-http \
  --env-vars-file env_updatedatastore.yaml 
```

## HOW-TO Use `updatedatastore` ##

### API ###
Example expected JSON to update Google Datastore with fresh mobile device data:

```json
{
	"key": "0123456789"
}
```

Example command line using `curl` and the above JSON to update Google Datastore with fresh movbile device data:

```
$ curl -X POST -d '{"key": "0123456789"}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/UpdateDatastore
```

It is recommended that Google [Cloud Scheduler](https://cloud.google.com/scheduler/) be used to schedule automatic calls to `updatedatastore`, for example: 

```
$ gcloud scheduler jobs list
ID                   LOCATION     SCHEDULE (TZ)                      TARGET_TYPE  STATE
UpdateDatastore      us-central1  */1 * * * * (America/Los_Angeles)  HTTP         ENABLED
```

### mdmtool ###
Fresh mobile device data can be downloaded and updated in Google Datastore by running `mdmtool`'s `updatedb` command:

```
$ mdmtool updatedb
Updating Datastore...  done.
updatedatastore Success
```

