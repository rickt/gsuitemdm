package slackdirectory

//
// GSuiteMDM slackdirectory Cloud Function
//

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	appname      string = os.Getenv("APPNAME")
	sm_config_id string = os.Getenv("SM_CONFIG_ID")
	sm_slacktoken_id string = os.Getenv("SM_SLACKTOKEN_ID")
)

// Search Google Datastore for a mobile device owner and return the associated phone number to Slack
func SlackDirectory(w http.ResponseWriter, r *http.Request) {
	var err error
	var devices []*gsuitemdm.DatastoreMobileDevice
	var l *logging.Client
	var text, token, user string

	// Null message body?
	if r.Body == nil {
		http.Error(w, "Error: Null message body", 400)
		return
	}

	// Not null, lets decode the x-www-form-urlencoded message body
	err = r.ParseForm()
	if err != nil {
		log.Printf("Error decoding Slack x-www-form-urlencoded message body: %s", err)
		http.Error(w, "Error decoding Slack x-www-form-urlencoded message body", 400)
		return
	}

	// Extract the pieces of the request we care about
	text = r.Form.Get("text")
	token = r.Form.Get("token")
	user = r.Form.Get("user_name")
	
  // Get a context
	ctx := context.Background()

	// Get the Slack token from Secret Manager
	slacktoken, err := gsuitemdm.GetSecret(ctx, sm_slacktoken_id)
	if err != nil {
		log.Printf("Error retrieving Slack token from Secret Manager", err)
		http.Error(w, "Error retrieving Slack token from Secret Manager", 400)
		return
	}

	// Check the key
	if token != slacktoken {
		log.Printf("Error: incorrect token sent with request")
		http.Error(w, "Not authorized, incorrect token", 401)
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

	// Initialise Stackdriver logging for this GCP project
	l, err = logging.NewClient(ctx, gs.C.ProjectID)

	// Register a Stackdriver logger instance for this app
	sl := l.Logger(appname)

	// Make sure query is not zero length
	if len(text) < 1 {
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

	// Range through the list of devices and search for a name using the text sent in the request from Slack
	for k := range devices {
		if strings.Contains(strings.ToUpper(devices[k].Name), strings.ToUpper(text)) {
			// Only return data if PhoneNumber exists
			if devices[k].PhoneNumber != "" {
				var p gsuitemdm.DirectoryData
				p.Name = devices[k].Name
				p.Email = devices[k].Email
				p.PhoneNumber = "(" + devices[k].PhoneNumber[0:3] + ") " + devices[k].PhoneNumber[3:6] + "-" + devices[k].PhoneNumber[6:10]
				dirdata = append(dirdata, p)
			}
		}
	}

	// Do we have any data to return? If so, return it
	if len(dirdata) > 0 {
		// Sort our data first
		sort.Sort(gsuitemdm.AllDirectoryData{dirdata})

		// We have valid search data to return
		var s string
		s = fmt.Sprintf("Users matching \"%s\": (%d)\n", text, len(dirdata))
		for n := range dirdata {
			s = s + fmt.Sprintf("%s: :dir_phone: %s :dir_email: `%s`\n", dirdata[n].Name, dirdata[n].PhoneNumber, dirdata[n].Email)
		}
		// Write the data
		w.Write([]byte(s))
		// Write a log entry
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: appname + " Success: " + strconv.Itoa(len(dirdata)) + " results returned for user=@" + user + ", q=" + text + " RemoteIP=" + gsuitemdm.GetIP(r)})
		return
	} else {
		// No data to return, say sorry
		var s string
		s = fmt.Sprintf("I'm sorry, but I was not able to find a user or group using your search term \"%s\"! :confused:\n", text)
		w.Write([]byte(s))
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: appname + " Success: 0 results returned for user=@" + user + ", q=" + text + " RemoteIP=" + gsuitemdm.GetIP(r)})
		return
	}
}

// EOF
