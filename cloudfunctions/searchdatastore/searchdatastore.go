package searchdatastore

//
// GSuiteMDM searchdatastore Cloud Function
//

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
	"strconv"
	"strings"
)

var (
	appname      string = os.Getenv("APPNAME")
	sm_apikey_id string = os.Getenv("SM_APIKEY_ID")
	sm_config_id string = os.Getenv("SM_CONFIG_ID")
)

// Search Google Datastore for a mobile devie
func SearchDatastore(w http.ResponseWriter, r *http.Request) {
	var err error
	var devices []*gsuitemdm.DatastoreMobileDevice
	var l *logging.Client
	var request gsuitemdm.SearchRequest

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

	// Ok, lets go deeper and check the message body. Was qtype= specified, and is it zero length?
	if len(request.QType) < 1 {
		log.Printf("Error: Query type not specified")
		http.Error(w, "Error: Query type not specified", 400)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Query type not specified"})
		return
	}

	// Do we support the specified query type
	if request.QType != "all" && request.QType != "email" && request.QType != "imei" && request.QType != "name" &&
		request.QType != "notes" && request.QType != "phone" && request.QType != "sn" && request.QType != "status" {
		log.Printf("Error: Invalid query type specified")
		http.Error(w, "Error: Invalid query type specified", 400)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Invalid query type specified"})
		return
	}

	// Query type is valid, lets check if the query string (q=) is not zero length. Only do this
	// if the query type is not 'all' as no 'q' parameter required if qtype==all
	if request.QType != "all" {
		// Check 'q=' since this is not a 'qtype=all' scenario
		if len(request.Q) < 1 {
			log.Printf("Error: Query search data cannot be zero length")
			http.Error(w, "Error: Query search data cannot be zero length", 400)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Query search data cannot be zero length"})
			return
		}
	}

	// Query type is valid and query string (q=) is not zero length, lets get the Datastore data
	dc, err := datastore.NewClient(ctx, gs.C.ProjectID)
	if err != nil {
		log.Printf("Error creating Datastore client: %s", err)
		http.Error(w, fmt.Sprintf("Error creating Datastore client: %s", err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error creating Datastore client: " + err.Error()})
		return
	}

	// Is this a domain-specific search?
	if request.Domain != "" && gs.IsDomainConfigured(request.Domain) == false {
		// Domain specified is invalid
		log.Printf("Invalid domain specified")
		http.Error(w, "Invalid domain specified", 400)
		return
	}

	// What kind of Datastore query do we make?
	if request.Domain != "" {
		// Perform a domain-specific search using a Datastore filter
		_, err = dc.GetAll(ctx, datastore.NewQuery(gs.C.DSNamekey).
			Filter("Domain =", request.Domain).
			Order(gs.C.DatastoreQueryOrderBy),
			&devices)
	} else {
		// Perform a full Datastore search with no filter
		_, err = dc.GetAll(ctx, datastore.NewQuery(gs.C.DSNamekey).
			Order(gs.C.DatastoreQueryOrderBy),
			&devices)
	}
	if err != nil {
		log.Printf("Error querying Datastore for all devices: %s", err)
		http.Error(w, fmt.Sprintf("Error querying Datastore for all devices: %s", err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error querying Datastore for all devices: " + err.Error()})
		return
	}

	// Return data right away if specified query type is "all"
	if request.QType == "all" {
		// Return some nice JSON data
		js, err := json.MarshalIndent(devices, "", "   ")
		if err != nil {
			log.Printf("Error marshaling JSON: %s", err)
			http.Error(w, fmt.Sprintf("Error marshaling JSON: %s", err), 500)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error marshaling JSON: " + err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		// Write a log entry
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: appname + " Success: " + strconv.Itoa(len(devices)) + " results returned"})
		return
	}

	// For all other query types, we much search through the device data
	var searchdata []*gsuitemdm.DatastoreMobileDevice

	for k := range devices {

		switch request.QType {
		case "email":
			if devices[k].Email == request.Q {
				searchdata = append(searchdata, devices[k])
				break
			}

		case "imei":
			if devices[k].IMEI == request.Q {
				searchdata = append(searchdata, devices[k])
				break
			}

		case "name":
			if strings.Contains(strings.ToUpper(devices[k].Name), strings.ToUpper(request.Q)) {
				searchdata = append(searchdata, devices[k])
				break
			}

		case "notes":
			if strings.Contains(strings.ToUpper(devices[k].Notes), strings.ToUpper(request.Q)) {
				searchdata = append(searchdata, devices[k])
				break
			}

		case "phone":
			if devices[k].PhoneNumber == request.Q {
				searchdata = append(searchdata, devices[k])
				break
			}

		case "sn":
			if devices[k].SN == request.Q {
				searchdata = append(searchdata, devices[k])
				break
			}

		case "status":
			if strings.Contains(strings.ToUpper(devices[k].Status), strings.ToUpper(request.Q)) {
				searchdata = append(searchdata, devices[k])
				break
			}

		default:
			http.Error(w, "Invalid query type specified", 400)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Invalid query type specified"})
			return
		}
	}

	// Do we have any data to return? If so, marshal into JSON and return it
	if len(searchdata) > 0 {
		// We have valid search data to return
		js, err := json.MarshalIndent(searchdata, "", "   ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling JSON: %s", err), 500)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error marshaling JSON: " + err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		// Write a log entry
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: appname + " Success: " + strconv.Itoa(len(searchdata)) + " results returned RemoteIP=" + gsuitemdm.GetIP(r)})
		return
	} else {
		// No data to return
		http.Error(w, "", 204)
		// Write a log entry
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: appname + " Success: 0 results returned RemoteIP=" + gsuitemdm.GetIP(r)})
		return
	}
}

// EOF
