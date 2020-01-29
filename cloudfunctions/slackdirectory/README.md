# gsuitemdm Cloud Function `slackdirectory` #

A [cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that can be used to search for a phone number using name or email as a search key. This cloud function is very similar to [`directory`](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/directory) but is specifically designed to format output as expected by Slack, and as such `slackdirectory` should be used as an HTTP backend for a Slack `/slash` command such as `/dir` or `/phone` (etc).

## HOW-TO Configure `slackdirectory` ##
`slackdirectory` uses a `.yaml` file containing several environment variables the cloud function reads during app startup. These environment variables point the app to the shared master cloud function configuration and the expected token to be received in each request from  Slack that are stored as [Secret Manager secrets](https://cloud.google.com/secret-manager/docs/managing-secrets). An example `.yaml` file for `slackdirectory`:

```yaml
APPNAME: slackdirectory
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_SLACKTOKEN_ID: projects/12334567890/secrets/gsuitemdm_slacktoken
```

## HOW-TO Deploy `slackdirectory` ##
```
$ gcloud functions deploy SlackDirectory --runtime go111 --trigger-http \
  --env-vars-file env_slackdirectory.yaml
```

## HOW-TO Use `slackdirectory` ##

### API ###

Example test command line that searches for phone numbers including the name "doe" and returns the output correctly formatted for display in Slack:

```
$ curl -X POST -d "token=0123456789&text=doe" \ 
  https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SlackDirectory
Users matching "doe": (2)
Jane Doe: :dir_phone: (213) 555-1313 :dir_email: `janed@foo.com`
John Doe: :dir_phone: (213) 555-1212 :dir_email: `johnd@foo.com`
```

