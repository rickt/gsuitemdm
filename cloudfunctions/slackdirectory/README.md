# gsuitemdm Cloud Function `directory` #

A [Cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that can be used to search for a phone number using email or name as a search key. 

## HOW-TO Deploy ##
`$ gcloud functions deploy Directory --runtime go111 --trigger-http --env-vars-file env.yaml`

## API Examples ##

Example test command line that searches for phone numbers including the name "doe":

```json
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

Example test command line that searches for phone numbers associated with the email address "johnd@foo.com":

```json
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

## TODO ##
* Add "all" query type to return all phone numbers
* Add option to return Slack-formatted data

