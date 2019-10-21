# gsuitemdm Cloud Function `updatesheet` #

A [Cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that updates a Google Sheet with fresh mobile device data from Datastore.

## HOW-TO Deploy ##
`$ gcloud functions deploy UpdateSheet --runtime go111 --trigger-http --env-vars-file env.yaml`


