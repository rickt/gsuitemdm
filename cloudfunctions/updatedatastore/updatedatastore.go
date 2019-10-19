package updatedatastore

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"log"
	"net/http"
	"os"
)

var (
	appname    string = os.Getenv("APPNAME")
	configfile string = os.Getenv("CONFIGFILE")
	key        string = os.Getenv("KEY")
)

func UpdateDatastore(w http.ResponseWriter, r *http.Request) {
	var domain string
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

	// Get existing Datastore data
	err = gs.GetDatastoreData()
	if err != nil {
		sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error getting existing Datastore data: " + fmt.Sprintf("%s", err)})
		return
	}

	// Range through the slice of configured domains
	for _, dm := range gs.C.Domains {
		domain = dm.DomainName

		// Get data about this domain's devices from the Admin SDK
		err = gs.GetAdminSDKDevices(domain)
		if err != nil {
			sl.Log(logging.Entry{Severity: logging.Error, Payload: "UpdateDatastore(): Error getting Admin SDK data for " + domain})
			return
		}

		// Range through this domain's list of devices and update it in Datastore
		for _, device := range gs.SDKData.Mobiledevices {
			err = gs.UpdateDatastoreDevice(device)
			if err != nil {
				sl.Log(logging.Entry{Severity: logging.Error, Payload: "Error converting device: " + fmt.Sprintf("%s", err)})
			}
		}
	}

	// Finished
	fmt.Fprintf(w, "Success\n")

	return
}
