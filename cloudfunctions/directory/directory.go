package directory

//
// GSuiteMDM directory Cloud Function
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
	"strings"
)

var (
	appname      string = os.Getenv("APPNAME")
	sm_apikey_id string = os.Getenv("SM_APIKEY_ID")
	sm_config_id string = os.Getenv("SM_CONFIG_ID")
)

// Search Google Datastore for a mobile device owner and return the associated phone number
func Directory(w http.ResponseWriter, r *http.Request) {
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
	if request.Key != strings.TrimSuffix(apikey, "\n") {
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

	// Do we support the specified query type? Directory supports only "email" and "name"
	if request.QType != "email" && request.QType != "name" {
	  log.Printf("Error: Invalid query type specified")
		http.Error(w, "Error: Invalid query type specified", 400)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Invalid query type specified"})
		return
	}

	// Query type is valid, lets check if the query string (q=) is not zero length
	if len(request.Q) < 1 {
	  log.Printf("Error: Query search data cannot be zero length")
		http.Error(w, "Error: Query search data cannot be zero length", 400)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Query search data cannot be zero length"})
		return
	}

	// Query type is valid and query string (q=) is not zero length, lets get the Datastore data
	dc, err := datastore.NewClient(ctx, gs.C.ProjectID)
	if err != nil {
	  log.Printf("Error creating Datastore client: %s", err)
		http.Error(w, fmt.Sprintf("Error creating Datastore client: %s", err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error creating Datastore client: " + err.Error()})
		return
	}

	// Perform a full Datastore search with no filter
	_, err = dc.GetAll(ctx, datastore.NewQuery(gs.C.DSNamekey).
		Order(gs.C.DatastoreQueryOrderBy),
		&devices)
	if err != nil {
	  log.Printf("Error querying Datastore for all devices: %s", err)
		http.Error(w, fmt.Sprintf("Error querying Datastore for all devices: %s", err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error querying Datastore for all devices: " + err.Error()})
		return
	}

	// Search for directory entries using the search type specified
	var dirdata []gsuitemdm.DirectoryData

	// Range through the list of devices and perform the correct type of search
	for k := range devices {

		switch request.QType {
		// Email search
		case "email":
			if devices[k].Email == request.Q {
				// Only return data if PhoneNumber exists
				if devices[k].PhoneNumber != "" {
					var p gsuitemdm.DirectoryData
					p.Name = devices[k].Name
					p.Email = devices[k].Email
					p.PhoneNumber = "(" + devices[k].PhoneNumber[0:3] + ") " + devices[k].PhoneNumber[3:6] + "-" + devices[k].PhoneNumber[6:10]
					dirdata = append(dirdata, p)
					break
				}
			}

		// Name search
		case "name":
			if strings.Contains(strings.ToUpper(devices[k].Name), strings.ToUpper(request.Q)) {
				// Only return data if PhoneNumber exists
				if devices[k].PhoneNumber != "" {
					var p gsuitemdm.DirectoryData
					p.Name = devices[k].Name
					p.Email = devices[k].Email
					p.PhoneNumber = "(" + devices[k].PhoneNumber[0:3] + ") " + devices[k].PhoneNumber[3:6] + "-" + devices[k].PhoneNumber[6:10]
					dirdata = append(dirdata, p)
					break
				}
			}

		default:
		  log.Printf("Error: Invalid query type specified")
			http.Error(w, "Error: Invalid query type specified", 400)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Invalid query type specified"})
			return
		}
	}

	// Do we have any data to return? If so, marshal into JSON and return it
	if len(dirdata) > 0 {
		// We have valid search data to return
		js, err := json.MarshalIndent(dirdata, "", "   ")
		if err != nil {
		  log.Printf("Error marshaling JSON: %s", err)
			http.Error(w, fmt.Sprintf("Error marshaling JSON: %s", err), 500)
			sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error marshaling JSON: " + err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	} else {
		// No data to return
		http.Error(w, "", 204)
		return
	}
}

// EOF
