package searchdatastore

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"log"
	"net/http"
	"os"
	// "strconv"
)

var (
	appname    string = os.Getenv("APPNAME")
	configfile string = os.Getenv("CONFIGFILE")
	domain     string = os.Getenv("DOMAIN")
	key        string = os.Getenv("KEY")
)

func SearchDatastore(w http.ResponseWriter, r *http.Request) {
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

	// Check URL parameters. Was qtype= specified, and is it zero length
	qt, ok := r.URL.Query()["qtype"]
	if !ok || len(qt[0]) < 1 {
		http.Error(w, "Query type not specified", 400)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Query type not specified"})
		return
	}

	// Do we support the specified query type
	if qt[0] != "all" && qt[0] != "email" && qt[0] != "imei" && qt[0] != "name" && qt[0] != "notes" &&
		qt[0] != "phone" && qt[0] != "sn" && qt[0] != "status" {
		http.Error(w, "Invalid query type specified", 400)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Invalid query type specified"})
		return
	}

	// Check specified query type
	switch qt[0] {
	// all
	case "all":
		var devices []gsuitemdm.DatastoreMobileDevice

		// Create a Datastore client
		dc, err := datastore.NewClient(ctx, gs.C.ProjectID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating Datastore client: %s", err), 500)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error creating Datastore client: " + err.Error()})
			return
		}

		// Build the Datastore query & get the list of devices
		_, err = dc.GetAll(ctx, datastore.NewQuery("MobileDevice").
			Order("Name"),
			&devices)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying Datastore for all devices: %s", err), 500)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error querying Datastore for all devices: " + err.Error()})
			return
		}

		// Return some nice data
		js, err := json.Marshal(devices)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling JSON: %s", err), 500)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error marshaling JSON: " + err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return

	case "email":
		fmt.Fprintf(w, "qtype=email\n")
		return

	case "imei":
		fmt.Fprintf(w, "qtype=imei\n")
		return

	case "name":
		fmt.Fprintf(w, "qtype=name\n")
		return

	case "notes":
		fmt.Fprintf(w, "qtype=notes\n")
		return

	case "phone":
		fmt.Fprintf(w, "qtype=phone\n")
		return

	case "sn":
		fmt.Fprintf(w, "qtype=sn\n")
		return

	case "status":
		fmt.Fprintf(w, "qtype=status\n")
		return

	default:
		http.Error(w, "Invalid query type specified", 400)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Invalid query type specified"})
		return

	}

	// Nearly finished
	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "gsuitemdm cloudfunction " + appname + " ended"})
		fmt.Fprintf(w, "gsuitemdm cloudfunction %s ended\n", appname)
	}

	// Finished
	fmt.Fprintf(w, "Success\n")

	return
}
