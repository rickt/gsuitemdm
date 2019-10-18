package searchdatastore

import (
	"cloud.google.com/go/logging"
	"context"
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

	if gs.C.Debug {
		sl.Log(logging.Entry{Severity: logging.Notice, Payload: "gsuitemdm cloudfunction " + appname + " ended"})
		fmt.Fprintf(w, "gsuitemdm cloudfunction %s ended\n", appname)
	}

	// Finished
	fmt.Fprintf(w, "Success\n")

	return
}
