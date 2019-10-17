package sheet2datastore

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	appname    string = os.Getenv("APPNAME")
	configfile string = os.Getenv("CONFIGFILE")
	domain     string = os.Getenv("DOMAIN")
)

func Sheet2Datastore(w http.ResponseWriter, r *http.Request) {
	var err error
	var l *logging.Client

	// Get a context
	ctx := context.Background()

	// Get a G Suite MDM Service
	gs, err := gsuitemdm.New(ctx, configfile)
	if err != nil {
		// Log to stderr, will be captured as a basic Stackdriver log
		log.Printf("Error: gsuitemdm cloudfunction %s could not start: %s", err)
		return
	}

	// Initialise Stackdriver logging for this GCP project
	l, err = logging.NewClient(ctx, gs.C.ProjectID)

	// Register a Stackdriver logger instance for this app
	sl := l.Logger(appname)

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "gsuitemdm cloudfunction " + appname + " started"})
	}

	// Get Admin SDK data
	err = gs.GetAdminSDKDevices(domain)
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error getting device data from G Suite Admin SDK: " + fmt.Sprintf("%s", err)})
		return
	}

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "Admin SDK reports " + strconv.Itoa(len(gs.SDKData.Mobiledevices)) + " mobile devices for domain " + domain})
	}

	// Get sheet data
	err = gs.GetSheetData()
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error retrieving Google Sheet data: " + fmt.Sprintf("%s", err)})
		return
	}

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "Google Sheet reports " + strconv.Itoa(len(gs.SheetData)-1) + " rows of data"})
	}

	// Get datastore data
	err = gs.GetDatastoreData()
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error retrieving Google Datastore data: " + fmt.Sprintf("%s", err)})
		return
	}

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "Google Datastore reports " + strconv.Itoa(len(gs.DatastoreData)) + " mobile devices for domain " + domain})
	}

	// Merge the data
	md := gs.MergeDatastoreAndSheetData()

	// Update the sheet
	err = gs.UpdateSheet(md)
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error updating Google Sheet: " + fmt.Sprintf("%s", err)})
		return
	}

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "Updated Google Sheet"})
	}

	// Update datastore
	count, err := gs.UpdateAllDatastoreData(domain)
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error updating Google Datastore: " + fmt.Sprintf("%s", err)})
		return
	}

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "Updated Google Datastore with " + strconv.Itoa(count) + " mobile devices for domain " + domain})
	}

	// Finished
	fmt.Fprintf(w, "SUCCESS\n")

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "gsuitemdm cloudfunction " + appname + " ended"})
	}
	return
}
