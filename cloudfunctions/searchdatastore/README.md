# gsuitemdm Cloud Function `searchdatastore` #

A [cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that searches for a mobile device in Google Datastore. Devices can be searched using owner name or email address, IMEI/SN/phone number or device status. 

## HOW-TO Configure `searchdatastore` ##
`searchdatastore` uses a `.yaml` file containing several environment variables the cloud function reads during app startup. These environment variables point the app to the shared master cloud function configuration and API key that are stored as [Secret Manager secrets](https://cloud.google.com/secret-manager/docs/managing-secrets). An example `.yaml` file for `searchdatastore`:

```yaml
APPNAME: searchdatastore
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```

## HOW-TO Deploy `searchdatastore` ##
```
$ gcloud functions deploy SearchDatastore \
  --runtime go111 \
  --trigger-http \
  --env-vars-file env_searchdatastore.yaml
```

## HOW-TO Use `searchdatastore` ##

### API ###
Example expected JSON to search for devices owned by 'john' (case insensitive owner name search):

```json
{
	"key": "0123456789",
	"q": "john",
	"qtype": "name"
}
```

Example expected JSON to search for devices that are currently blocked in the domain 'foo.com':
```json
{
	"domain": "foo.com",
	"key": "0123456789",
	"q": "BLOCKED",
	"qtype": "status"
}
```

Example expected JSON to search for a device with an IMEI of 111111111111111:
```json
{
	"key": "0123456789",
	"q": "111111111111111",
	"qtype": "imei"
}
```

Example command line using `curl` to search for devices owned by 'john' (case insensitive owner name search):

```
$ curl -X POST -d '{"key": "0123456789", "qtype": "name", "q": "john"}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastore
[
   {
      "Color": "black",
      "CompromisedStatus": "No compromise detected",
      "Domain": "foo.com",
      "DeveloperMode": false,
      "Email": "johnd@foo.com",
      "IMEI": "111111111111111",
      "Model": "iPhone 11 Pro",
      "Name": "John Doe",
      "Notes": "this is John's 3rd phone",
      "OS": "iOS 13.2.1",
      "OSBuild": "16F203",
      "PhoneNumber": "2135551212",
      "RAM": "64",
      "SN": "Z01ABCD0ABCD",
      "Status": "APPROVED",
      "SyncFirst": "2019-01-18T17:13:54.297Z",
      "SyncLast": "2019-10-21T16:17:02.935Z",
      "Type": "IOS_SYNC",
      "UnknownSources": false,
      "USBADB": false,
      "WifiMac": "aa:11:bb:22:cc:33"
   }
```

By way of illustration, the above same data would be returned with the following searches:

```
$ curl -X POST -d '{"key": "0123456789", "qtype": "email", "q": "johnd@foo.com"}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastore

$ curl -X POST -d '{"key": "0123456789", "qtype": "imei", "q": "111111111111111"}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastore

$ curl -X POST -d '{"key": "0123456789", "qtype": "phone", "q": "2135551212"}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastore

$ curl -X POST -d '{"key": "0123456789", "qtype": "sn", "q": "Z01ABCD0ABCD"}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastore
```

### `mdmtool` ###
Similarly, `mdmtool` can be used in many ways to search for devices. Some `mdmtool` examples to return the same data as shown above:

```
$ mdmtool search -n john

$ mdmtool search -e johnd@foo.com

$ mdmtool search -t BLOCKED -d foo.com

$ mdmtool search -i 012345678901234

$ mdmtool search -p 2135551212

$ mdmtool search -s Z01ABCD0ABCD
```

