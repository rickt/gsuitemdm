#! /bin/bash

# change this to point to your own GCP project
PROJECT="mdm-updater"

CLOUDFUNCTIONS="approvedevice blockdevice deletedevice directory searchdatastore slackdirectory updatedatastore updatesheet wipedevice"

for FUNCTION in $CLOUDFUNCTIONS
do
	echo "*** $FUNCTION ***"
	cd $FUNCTION
	go get -u github.com/rickt/gsuitemdm
	go clean && go build
	gcloud config set project $PROJECT
	gcloud functions deploy $FUNCTION --runtime go111 --trigger-http --env-vars-file env_${FUNCTION}.yaml
done

# EOF
