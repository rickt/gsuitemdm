package updatesheet

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"log"
	"net/http"
	"os"
)

// $ gcloud functions deploy UpdateSheet --runtime go111 --trigger-http --env-vars-file env.yaml

var (
	appname    string = os.Getenv("APPNAME")
	configfile string = os.Getenv("CONFIGFILE")
	key        string = os.Getenv("KEY")
)

func UpdateSheet(w http.ResponseWriter, r *http.Request) {
	var err error
	var l *logging.Client

	// Has the correct key been sent with the request?
	sk, ok := r.URL.Query()["key"]
	if !ok || len(sk[0]) < 1 || sk[0] != key {
		log.Printf("Error: incorrect key sent with request: %s", err)
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
	d := r.URL.Query().Get("debug")
	if len(d) != 0 {
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
