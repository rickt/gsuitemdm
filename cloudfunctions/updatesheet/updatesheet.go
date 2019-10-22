package updatesheet

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

var (
	appname    string = os.Getenv("APPNAME")
	configfile string = os.Getenv("CONFIGFILE")
	key        string = os.Getenv("KEY")
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

	// Get Google Sheet data
	err = gs.GetSheetData()
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error retrieving Google Sheet data: " + fmt.Sprintf("%s", err)})
		return
	}

	// Get existing Datastore data
	err = gs.GetDatastoreData()
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error retrieving Google Datastore data: " + fmt.Sprintf("%s", err)})
		return
	}

	// Merge the data
	var md []gsuitemdm.DatastoreMobileDevice
	md = gs.MergeDatastoreAndSheetData()

	// Update the Google Sheet
	err = gs.UpdateSheet(md)
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error updating Google Sheet: " + fmt.Sprintf("%s", err)})
		return
	}

	// Finished
	sl.Log(logging.Entry{Severity: logging.Notice, Payload: appname + " Success"})
	fmt.Fprintf(w, "%s Success\n", appname)

	return
}

// EOF
