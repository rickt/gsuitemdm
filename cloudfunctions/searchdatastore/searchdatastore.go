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
	"strings"
)

var (
	appname    string = os.Getenv("APPNAME")
	configfile string = os.Getenv("CONFIGFILE")
	key        string = os.Getenv("KEY")
)

// $ gcloud functions deploy SearchDatastore --runtime go111 --trigger-http --env-vars-file env.yaml

func SearchDatastore(w http.ResponseWriter, r *http.Request) {
	var err error
	var devices []*gsuitemdm.DatastoreMobileDevice
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
	dbg := r.URL.Query().Get("debug")
	if len(dbg) != 0 {
		gs.C.Debug = true
	}

	// Initialise Stackdriver logging for this GCP project
	l, err = logging.NewClient(ctx, gs.C.ProjectID)

	// Register a Stackdriver logger instance for this app
	sl := l.Logger(appname)

	// Check URL parameters. Was qtype= specified, and is it zero length
	var qt []string
	var qtype string

	qt, ok = r.URL.Query()["qtype"]
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

	// Query type is valid, continue
	qtype = qt[0]

	// Query type is valid, lets check if the query string (q=) is not zero length. Only do this
	// if the query type is not 'all' as no 'q' parameter required if qtype==all
	var qs []string
	var qstring string

	if qt[0] != "all" {
		// Check 'q=' since this is not a 'qtype=all' scenario
		qs, ok = r.URL.Query()["q"]
		if !ok || len(qs[0]) < 1 {
			http.Error(w, "Query search data cannot be zero length", 400)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Query search data cannot be zero length"})
			return
		} else {
			// Query string is valid, continue
			qstring = qs[0]
		}
	}

	// Query type is valid and query string (q=) is not zero length, lets get the Datastore data
	// Create a Datastore client
	dc, err := datastore.NewClient(ctx, gs.C.ProjectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Datastore client: %s", err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error creating Datastore client: " + err.Error()})
		return
	}

	// Was domain URL parameter specified?
	var domain string

	// Check if domain= was specified
	d, ok := r.URL.Query()["domain"]
	// Yes, domain= was specified, but is the domain valid?
	if ok {
		if len(d[0]) < 1 || gs.IsDomainConfigured(d[0]) == false {
			// Domain specified is invalid
			log.Printf("Invalid domain specified")
			http.Error(w, "Invalid domain specified", 200)
			return
		} else {
			// Domain is valid, continue
			domain = d[0]
		}
	}

	// What kind of Datastore query do we make?
	if len(domain) > 1 {
		// Perform a domain-specific search using a Datastore filter
		_, err = dc.GetAll(ctx, datastore.NewQuery("MobileDevice").
			Filter("Domain =", domain).
			Order(gs.C.DatastoreQueryOrderBy),
			&devices)
	} else {
		// Perform a full Datastore search with no filter
		_, err = dc.GetAll(ctx, datastore.NewQuery("MobileDevice").
			Order(gs.C.DatastoreQueryOrderBy),
			&devices)
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying Datastore for all devices: %s", err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error querying Datastore for all devices: " + err.Error()})
		return
	}

	// Return data right away if specified query type is "all"
	if qtype == "all" {
		// Return some nice JSON data
		js, err := json.MarshalIndent(devices, "", "   ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling JSON: %s", err), 500)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error marshaling JSON: " + err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}

	// For all other query types, we much search through the device data
	var searchdata []*gsuitemdm.DatastoreMobileDevice

	for k := range devices {

		switch qtype {
		case "email":
			if devices[k].Email == qstring {
				searchdata = append(searchdata, devices[k])
			}

		case "imei":
			if devices[k].IMEI == qstring {
				searchdata = append(searchdata, devices[k])
			}

		case "name":
			if strings.Contains(strings.ToUpper(devices[k].Name), strings.ToUpper(qstring)) {
				searchdata = append(searchdata, devices[k])
			}

		case "notes":
			if strings.Contains(strings.ToUpper(devices[k].Notes), strings.ToUpper(qstring)) {
				searchdata = append(searchdata, devices[k])
			}

		case "phone":
			if devices[k].PhoneNumber == qstring {
				searchdata = append(searchdata, devices[k])
			}

		case "sn":
			if devices[k].SN == qstring {
				searchdata = append(searchdata, devices[k])
			}

		case "status":
			if strings.Contains(strings.ToUpper(devices[k].Status), strings.ToUpper(qstring)) {
				searchdata = append(searchdata, devices[k])
			}

		default:
			http.Error(w, "Invalid query type specified", 400)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Invalid query type specified"})
			return
		}
	}

	// Return some nice JSON data
	js, err := json.MarshalIndent(searchdata, "", "   ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling JSON: %s", err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error marshaling JSON: " + err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	return
}
