package approvedevice

//
// GSuiteMDM approvedevice Cloud Function
//

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rickt/gsuitemdm"
	admin "google.golang.org/api/admin/directory/v1"
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

// Approve a mobile device using the G Suite Admin SDK
func ApproveDevice(w http.ResponseWriter, r *http.Request) {
	var aa = &admin.MobileDeviceAction{
		Action: "approve",
	}
	var as *admin.Service
	var cid string
	var devices []*gsuitemdm.DatastoreMobileDevice
	var err error
	var l *logging.Client
	var request gsuitemdm.ActionRequest

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

	// Correct action specified?
	if request.Action != "approve" {
		log.Printf("Error: Invalid action specified")
		http.Error(w, "Invalid request (invalid action specified)", 400)
		return
	}

	// Check if the request is valid
	if (request.IMEI == "" && request.SN == "") || (request.IMEI != "" && request.SN != "") {
		log.Printf("Error: Invalid request (IMEI or SN not specified)")
		http.Error(w, "Invalid request (IMEI or SN not specified)", 400)
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

	// Was the (required) domain specified?
	if request.Domain == "" || gs.IsDomainConfigured(request.Domain) == false {
		// Domain specified is invalid
		log.Printf("Error: Invalid domain specified")
		http.Error(w, "Invalid domain specified", 400)
		return
	}

	// Ok, the action + domain are valid, lets get the Datastore data
	dc, err := datastore.NewClient(ctx, gs.C.ProjectID)
	if err != nil {
		log.Printf("Error creating Datastore client: %s", err)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error creating Datastore client: " + err.Error()})
		http.Error(w, fmt.Sprintf("Error creating Datastore client: %s", err), 500)
		return
	}

	// Get existing Datastore data for this domain
	_, err = dc.GetAll(ctx, datastore.NewQuery(gs.C.DSNamekey).
		Filter("Domain =", request.Domain).
		Order(gs.C.DatastoreQueryOrderBy),
		&devices)
	if err != nil {
		log.Printf("Error querying Datastore for devices in domain %s: %s", request.Domain, err)
		http.Error(w, fmt.Sprintf("Error querying Datastore for devices in domain %s: %s", request.Domain, err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error querying Datastore for devices in domain " + request.Domain + ": " + err.Error()})
		return
	}

	// Iterate through this domain's devices to find the specified device
	var found bool = false
	var device *gsuitemdm.DatastoreMobileDevice

	for _, d := range devices {
		// How are we identifying the mobile device? SN or IMEI?
		switch {
		// IMEI
		case request.IMEI != "":
			if strings.Replace(d.IMEI, " ", "", -1) == strings.Replace(request.IMEI, " ", "", -1) && d.Domain == request.Domain {
				found = true
				device = d
				break
			}
		// SN
		case request.SN != "":
			if strings.Replace(d.SN, " ", "", -1) == strings.Replace(request.SN, " ", "", -1) && d.Domain == request.Domain {
				found = true
				device = d
				break
			}
		}
	}

	// Did we find the specified device?
	if found != true {
		log.Printf("Error: Device not found")
		http.Error(w, "Error: Device not found", 400)
		return
	}

	// Check if device has the correct G Suite MDM status. Valid states for approve are:
	// PENDING and BLOCKED
	if device.Status != "PENDING" && device.Status != "BLOCKED" {
		log.Printf("Error: Device found but not in BLOCKED or PENDING states")
		http.Error(w, fmt.Sprintf("Error: Device found but not in BLOCKED or PENDING states (status=%s)", device.Status), 400)
		return
	}

	// Was `confirm: true` sent along with the request?
	if request.Confirm != true {
		log.Printf("Error: Device found and in BLOCKED or PENDING states but no CONFIRM sent")
		http.Error(w, "Error: Device found and in BLOCKED or PENDING states but no CONFIRM sent", 400)
		return
	}

	// Confirm was sent, lets approve the device. Get this domain's CustomerID first
	cid, err = gs.GetDomainCustomerID(request.Domain)
	if err != nil {
		log.Printf("Error getting CustomerID for domain %s: %s", request.Domain, err)
		http.Error(w, fmt.Sprintf("Error getting CustomerID for domain %s: %s", request.Domain, err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error getting CustomerID for domain " + request.Domain + ": " + err.Error()})
		return
	}

	// Authenticate with the Admin SDK for this domain
	as, err = gs.AuthenticateWithDomain(cid, request.Domain, gs.C.ActionScope)
	if err != nil {
		log.Printf("Error authenticating with the Admin SDK for domain %s: %s", request.Domain, err)
		http.Error(w, fmt.Sprintf("Error authenticating with the Admin SDK for domain %s: %s", request.Domain, err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error authenticating with the Admin SDK for domain " + request.Domain + ": " + err.Error()})
		return
	}

	// Approve the device
	err = as.Mobiledevices.Action(cid, device.ResourceId, aa).Do()
	if err != nil {
		log.Printf("Error approving device %s in domain %s: %s", device.ResourceId, request.Domain, err)
		http.Error(w, fmt.Sprintf("Error approving device %s in domain %s: %s", device.ResourceId, request.Domain, err), 500)
		sl.Log(logging.Entry{Severity: logging.Warning, Payload: "Error approving device " + device.ResourceId + " in domain " + request.Domain + ": " + err.Error()})
		return
	}

	// Finished, write a log entry
	log.Printf("")
	sl.Log(logging.Entry{Severity: logging.Notice, Payload: appname + " Success: SN=" + device.SN + " Owner=" + device.Email + " RemoteIP=" + GetIP(r)})
	fmt.Fprintf(w, "%s Success\n", appname)

	return
}

// EOF
