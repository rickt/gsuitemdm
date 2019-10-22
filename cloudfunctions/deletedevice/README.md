# gsuitemdm Cloud Function `deletedevice` #

A [Cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that deletes a mobile device using the [Admin SDK](https://developers.google.com/admin-sdk).

## HOW-TO Deploy ##
`$ gcloud functions deploy WipeDevice --runtime go111 --trigger-http --env-vars-file env.yaml`


