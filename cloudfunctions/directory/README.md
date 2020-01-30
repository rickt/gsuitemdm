# gsuitemdm Cloud Function `directory` #

A [cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package providing an API used to search for phone number data (name, phone number) among all tracked/configured mobile devices using email address, name or phone number as a search key. 

This API might be useful for orgs who want a handy way to search for phone numbers among their company mobile device-using staff. With `directory` you have an always up-to-date, automatically-updated mobile phone directory API that searches all device data in all configured G Suite domains. 

## HOW-TO Configure `directory` ##
`directory` uses a `.yaml` file containing several environment variables the cloud function reads during app startup. These environment variables point the app to the shared master cloud function configuration and API key that are stored as [Secret Manager secrets](https://cloud.google.com/secret-manager/docs/managing-secrets). An example `.yaml` file for `directory`:

```yaml
APPNAME: directory
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```

## HOW-TO Deploy `directory` ##
```
$ gcloud functions deploy Directory \
  --runtime go111 \
  --trigger-http \
  --env-vars-file env_directory.yaml
```

## HOW-TO Use `directory` ##

### API ###
Example expected JSON to search for phone numbers including the name "doe":

```json
{
	"key": "0123456789",
	"q": "doe",
	"qtype": "name"
}
```

Example expected JSON to search for phone numbers associated with the email address "johnd@foo.com":
```json
{
	"key": "0123456789",
	"q": "johnd@foo.com",
	"qtype": "email"
}
```

Example command line using `curl` and the above JSON that searches for phone numbers including the name "doe":

```
$ curl -X POST -d '{"key": "0123456789", "qtype": "name", "q": "doe"}' \ 
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/Directory
[
   {
      "name": "John Doe",
      "email": "johnd@foo.com",
      "phonenumbner": "(213) 555-1212"
   },
   {
      "name": "Jane Doe",
      "email": "janed@foo.com",
      "phonenumbner": "(213) 555-1313"
   }
]
```

Example command line using `curl` and the above JSON that searches for phone numbers associated with the email address "johnd@foo.com":

```
$ curl -X POST -d '{"key": "0123456789", "qtype": "email", "q": "johnd@foo.com"}' \
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/Directory
[
   {
      "name": "John Doe",
      "email": "johnd@foo.com",
      "phonenumbner": "(213) 555-1212"
   }
]
```

### `mdmtool` ###
Example command line using `mdmtool` to search for phone numbers including the name "doe":
```
$ mdmtool dir -n doe
----------------------+----------------+------------------------------------------
Name                  | Phone Number   | Email 
----------------------+----------------+------------------------------------------
Jane Doe              | (213) 555-1313 | janed@foo.com
John Doe              | (213) 555-1212 | johnd@foo.com
----------------------+----------------+------------------------------------------
Search returned 2 results.
```

Example command line using `mdmtool` to search for phone numbers associated with the email address "johnd@foo.com":
```
$ mdmtool dir -e johnd@foo.com
----------------------+----------------+------------------------------------------
Name                  | Phone Number   | Email 
----------------------+----------------+------------------------------------------
John Doe              | (213) 555-1212 | johnd@foo.com
----------------------+----------------+------------------------------------------
Search returned 1 results.
```
