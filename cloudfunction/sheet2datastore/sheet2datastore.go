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
	key        string = os.Getenv("KEY")
)

// Example deploy command line:
//
// $ cd gsuitemdm/demo/cloudfunction
// $ cp env_example.yaml env.yaml
// # edit env.yaml as appropriate to your environment
// $ cp cloudfunction_conf_example.json cloudfunction_conf.json
// # edit cloudfunction_conf.json as appropriate to your environment
// $ gcloud config set project <yourprojectname>
// $ gcloud functions deploy Sheet2Datastore --runtime go111 --trigger-http --env-vars-file env.yaml

func Sheet2Datastore(w http.ResponseWriter, r *http.Request) {
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

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "gsuitemdm cloudfunction " + appname + " started"})
		fmt.Fprintf(w, "gsuitemdm cloudfunction %s started\n", appname)
	}

	// Get Admin SDK data
	err = gs.GetAdminSDKDevices(domain)
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error getting device data from G Suite Admin SDK: " + fmt.Sprintf("%s", err)})
		return
	}

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "Admin SDK reports " + strconv.Itoa(len(gs.SDKData.Mobiledevices)) + " mobile devices for domain " + domain})
		fmt.Fprintf(w, "Admin SDK reports %d mobile devices for domain %s\n", len(gs.SDKData.Mobiledevices), domain)
	}

	// Get sheet data
	err = gs.GetSheetData()
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error retrieving Google Sheet data: " + fmt.Sprintf("%s", err)})
		return
	}

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "Google Sheet reports " + strconv.Itoa(len(gs.SheetData)-1) + " rows of data"})
		fmt.Fprintf(w, "Google Sheet reports %d mobile devices for domain %s\n", len(gs.SheetData)-1, domain)
	}

	// Get datastore data
	err = gs.GetDatastoreData()
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error retrieving Google Datastore data: " + fmt.Sprintf("%s", err)})
		return
	}

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "Google Datastore reports " + strconv.Itoa(len(gs.DatastoreData)) + " mobile devices for domain " + domain})
		fmt.Fprintf(w, "Google Datastore reports %d mobile devices for domain %s\n", len(gs.DatastoreData), domain)
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
		fmt.Fprintf(w, "Updated Google Sheet\n")
	}

	// Update datastore
	count, err := gs.UpdateAllDatastoreData(domain)
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error updating Google Datastore: " + fmt.Sprintf("%s", err)})
		return
	}

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "Updated Google Datastore with " + strconv.Itoa(count) + " mobile devices for domain " + domain})
		fmt.Fprintf(w, "Updated Google Datastore with %d mobile devices for domain %s\n", count, domain)
	}

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "gsuitemdm cloudfunction " + appname + " ended"})
		fmt.Fprintf(w, "gsuitemdm cloudfunction %s ended\n", appname)
	}

	// Finished
	fmt.Fprintf(w, "Success\n")

	return
}