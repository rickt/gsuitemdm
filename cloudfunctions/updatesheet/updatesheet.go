package updatesheet

//
// GSuiteMDM updatesheet Cloud Function
//

import (
	"cloud.google.com/go/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	appname      string = os.Getenv("APPNAME")
	sm_apikey_id string = os.Getenv("SM_APIKEY_ID")
	sm_config_id string = os.Getenv("SM_CONFIG_ID")
)

// Update the Google Sheet with fresh data from Google Datastore
func UpdateSheet(w http.ResponseWriter, r *http.Request) {
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

	// Get a context
	ctx := context.Background()

	// Get the API key from Secret Manager
	apikey, err := gsuitemdm.GetSecret(ctx, sm_apikey_id)
	if err != nil {
		log.Printf("Error retrieving API key from Secret Manager", err)
		http.Error(w, "Error retrieving API key from Secret Manager", 400)
		return
	}

	// Check that the API key sent with the request matches
	if request.Key != strings.TrimSpace(apikey) {
		log.Printf("Error: Incorrect key sent with request")
		http.Error(w, "Not authorized", 401)
		return
	}

	// Get our app configuration from Secret Manager
	config, err := gsuitemdm.GetSecret(ctx, sm_config_id)
	if err != nil {
		log.Printf("Error retrieving app configuration from Secret Manager: %s", err)
		http.Error(w, "Error retrieving app configuration from Secret Manager", 400)
		return
	}

	// Get a G Suite MDM Service
	gs, err := gsuitemdm.New(ctx, config)
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

	// Get Google Sheet data
	err = gs.GetSheetData()
	if err != nil {
	  log.Printf("Error retrieving Google Sheet data: %s", err)
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error retrieving Google Sheet data: " + fmt.Sprintf("%s", err)})
		return
	}

	// Get existing Datastore data
	err = gs.GetDatastoreData()
	if err != nil {
	  log.Printf("Error retrieving Google Datastore data: %s", err)
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error retrieving Google Datastore data: " + fmt.Sprintf("%s", err)})
		return
	}

	// Merge the data
	var md []gsuitemdm.DatastoreMobileDevice
	md = gs.MergeDatastoreAndSheetData()

	// Update the Google Sheet
	err = gs.UpdateSheet(md)
	if err != nil {
	  log.Printf("Error updating Google Sheet: %s", err)
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error updating Google Sheet: " + fmt.Sprintf("%s", err)})
		return
	}

	// Finished
	sl.Log(logging.Entry{Severity: logging.Notice, Payload: appname + " Success RemoteIP=" + gsuitemdm.GetIP(r)})
	fmt.Fprintf(w, "%s Success\n", appname)

	return
}

// EOF
