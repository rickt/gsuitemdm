# gsuitemdm Cloud Function `updatedatastore` #

A [Cloud Function](https://cloud.google.com/functions/) component of the [gsuitemdm](https://github.com/rickt/gsuitemdm) package that updates mobile device data in Google Datastore with the latest mobile device data from the [Admin SDK](https://developers.google.com/admin-sdk) as well as any local notes in the Google Sheet. 

## HOW-TO Deploy ##
`$ gcloud functions deploy UpdateDatastore --runtime go111 --trigger-http --env-vars-file env.yaml --memory 512MB`

