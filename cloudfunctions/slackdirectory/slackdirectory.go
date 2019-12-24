package slackdirectory

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
	slacktoken string = os.Getenv("SLACKTOKEN")
)

// Search Google Datastore for a mobile device owner and return the associated phone number
func Directory(w http.ResponseWriter, r *http.Request) {
	var err error
	var devices []*gsuitemdm.DatastoreMobileDevice
	var l *logging.Client
	var slackrequest gsuitemdm.SlackRequest

	// Null message body?
	if r.Body == nil {
		http.Error(w, "Error: Null message body", 400)
		return
	}

	// Not null, lets decode the message body
	err = json.NewDecoder(r.Body).Decode(&slackrequest)
	if err != nil {
		log.Printf("Error decoding Slack JSON message body: %s", err)
		http.Error(w, "Error decoding Slack JSON message body", 400)
		return
	}

	// Check the key
	if slackrequest.Token != slacktoken {
		log.Printf("Error: incorrect token sent with request")
		http.Error(w, "Not authorized, incorrect token", 401)
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

	// Initialise Stackdriver logging for this GCP project
	l, err = logging.NewClient(ctx, gs.C.ProjectID)

	// Register a Stackdriver logger instance for this app
	sl := l.Logger(appname)

	// Make sure query is not zero length
	if len(slackrequest.Text) < 1 {
		http.Error(w, "Query search data cannot be zero length", 400)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Query search data cannot be zero length"})
		return
	}

	// OK, lets get the Datastore data
	dc, err := datastore.NewClient(ctx, gs.C.ProjectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Datastore client: %s", err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error creating Datastore client: " + err.Error()})
		return
	}

	// Perform a full Datastore search with no filter
	_, err = dc.GetAll(ctx, datastore.NewQuery(gs.C.DSNamekey).
		Order(gs.C.DatastoreQueryOrderBy),
		&devices)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying Datastore for all devices: %s", err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error querying Datastore for all devices: " + err.Error()})
		return
	}

	// Search for directory entries using the search type specified
	var dirdata []gsuitemdm.DirectoryData

	// Range through the list of devices and perform the correct type of search
	for k := range devices {

		if strings.Contains(strings.ToUpper(devices[k].Name), strings.ToUpper(slackrequest.Text)) {
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
	}

	// Do we have any data to return? If so, marshal into JSON and return it
	if len(dirdata) > 0 {
		// We have valid search data to return
		js, err := json.MarshalIndent(dirdata, "", "   ")
		if err != nil {
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
