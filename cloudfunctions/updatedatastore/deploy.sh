#! /bin/bash

# change this to point to your own GCP project
PROJECT="mdm-updater"

go get -u github.com/rickt/gsuitemdm
go build
gcloud config set project $PROJECT
gcloud functions deploy UpdateDatastore --runtime go111 --trigger-http --env-vars-file env_updatedatastore.yaml --memory 512MB
