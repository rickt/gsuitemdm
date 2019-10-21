package updatedatastore

import (
	"cloud.google.com/go/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"log"
	"net/http"
	"os"
)

// Example deploy command line:
// $ gcloud functions deploy UpdateDatastore --runtime go111 --trigger-http --env-vars-file env.yaml

// Example command line to trigger a Google Datastore update:
// $ curl -X POST -d '{"key": "0123456789", "debug": false}' https://us-central1-<YOURGCPPROJECTNAME>.cloudfunctions.net/UpdateDatastore

var (
	appname    string = os.Getenv("APPNAME")
	configfile string = os.Getenv("CONFIGFILE")
	key        string = os.Getenv("KEY")
)

func UpdateDatastore(w http.ResponseWriter, r *http.Request) {
	var domain string
	var err error
	var l *logging.Client
	var request gsuitemdm.UpdateRequest

	// Null message body?
	if r.Body == nil {
		http.Error(w, "Error: Null message body", 400)
		return
	}

	// Not null, lets decode the message body
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error decoding JSON message body: %s", err)
		http.Error(w, "Error decoding JSON message body", 400)
		return
	}

	// Check the key
	if request.Key != key {
		log.Printf("Error: incorrect key sent with request")
		http.Error(w, "Not authorized", 401)
		return
	}

	// Get a context
	ctx := context.Background()

	// Get a G Suite MDM Service
	gs, err := gsuitemdm.New(ctx, configfile)
	if err != nil {
		// Log to stderr, will be captured as a basic Stackdriver log
		log.Printf("Error: gsuitemdm cloudfunction %s could not start: %s", err)
		return
	}

	// Debug mode?
	if request.Debug == true {
		gs.C.Debug = true
	}

	// Initialise Stackdriver logging for this GCP project
	l, err = logging.NewClient(ctx, gs.C.ProjectID)

	// Register a Stackdriver logger instance for this app
	sl := l.Logger(appname)

	// Get existing Datastore data
	err = gs.GetDatastoreData()
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error getting existing Datastore data: " + fmt.Sprintf("%s", err)})
		return
	}

	// Get Google Sheet data
	err = gs.GetSheetData()
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error getting Google Sheet data: " + fmt.Sprintf("%s", err)})
		return
	}

	// Range through the slice of configured domains
	for _, dm := range gs.C.Domains {
		domain = dm.DomainName

		// Get data about this domain's devices from the Admin SDK
		err = gs.GetAdminSDKDevices(domain)
		if err != nil {
			sl.Log(logging.Entry{Severity: logging.Error, Payload: "UpdateDatastore(): Error getting Admin SDK data for " + domain})
			return
		}

		// Range through this domain's list of devices and update it in Datastore
		for _, device := range gs.SDKData.Mobiledevices {
			err = gs.UpdateDatastoreDevice(device)
			if err != nil {
				sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error converting device: " + fmt.Sprintf("%s", err)})
			}
		}
	}

	// Finished
	sl.Log(logging.Entry{Severity: logging.Notice, Payload: appname + " Success"})
	fmt.Fprintf(w, "%s Success\n", appname)

	return
}

// EOF
