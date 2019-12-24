# gsuitemdm Cloud Function `slackdirectory` #

A [Cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that is used as HTTP backend for a Slack `/slash` command such as `/dir` or `/phone` (etc).

## HOW-TO Deploy ##
`$ gcloud functions deploy SlackDirectory --runtime go111 --trigger-http --env-vars-file env.yaml`

## API Examples ##

Example test command line that searches for phone numbers including the name "doe":

```
$ curl -d "token=0123456789&text=ick" -X POST https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/SlackDirectory
Users matching "ick": (3)
Rick Tait: :dir_phone: (213) 555-1212 :dir_email: `rickt@rickt.org`
Rick Tate: :dir_phone: (213) 555-1313 :dir_email: `rickt@foo.com`
Rick Tayt: :dir_phone: (213) 555-1414 :dir_email: `rickt@bar.com`
```

## TODO ##
