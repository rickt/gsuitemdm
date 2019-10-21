# gsuitemdm Cloud Function `searchdatastore` #

A [Cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that searches for a mobile device in Google Datastore. 

## HOW-TO Deploy ##
`$ gcloud functions deploy SearchDatastore --runtime go111 --trigger-http --env-vars-file env.yaml`

## API Examples ##
Example test command line that searches Datastore for devices owned by 'john' (case insensitive owner name search):

`
$ curl -X POST -d '{"key": "0123456789", "qtype": "name", "q": "rick"}' https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastorePost
[
   {
      "Color": "black",
      "CompromisedStatus": "No compromise detected",
      "Domain": "foo.com",
      "DeveloperMode": false,
      "Email": "johnd@foo.com",
      "IMEI": "012345678901234",
      "Model": "iPhone XR",
      "Name": "John Doe",
      "Notes": "this is John's 3rd phone",
      "OS": "iOS 12.3.1",
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
`

By way of illustration, the above same data would be returned with the following searches:

`
$ curl -X POST -d '{"key": "0123456789", "qtype": "email", "q": "johnd@foo.com"}' https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastorePost
$ curl -X POST -d '{"key": "0123456789", "qtype": "imei", "q": "012345678901234"}' https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastorePost
$ curl -X POST -d '{"key": "0123456789", "qtype": "notes", "q": "3rd"}' https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastorePost
$ curl -X POST -d '{"key": "0123456789", "qtype": "phone", "q": "2135551212"}' https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastorePost
$ curl -X POST -d '{"key": "0123456789", "qtype": "sn", "q": "Z01ABCD0ABCD"}' https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SearchDatastorePost
`

