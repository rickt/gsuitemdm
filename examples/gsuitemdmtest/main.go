package main

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"log"
	"os"
	"strconv"
)

// Sample code showing how to use the gsuitemdm package.
//
//	What it does:
//		* downloads all mobile device data for a G Suite domain's MDM-managed devices using the G Suite Admin SDK
//		* downloads all mobile device data from a tracking Google sheet
//		* downloads all mobile device data from Google Datastore
//		* merges all the data
//		* updates the Google sheet
//		* updates Google Datastore
//
// Basic requirements:
// 	1) a G Suite domain with MDM enabled and > 1 mobile device configured w/MDM
//		2) a properly-setup main configuration file. be sure that 'sheetcreds', 'sheetid',
//		   'sheetwho' and 'projectid' are configured correctly! an example configuration
//			file is included (gsuitemdmtest_example_conf.json)
//		3) TODO
//
//	Instructions:
//		1) add your own app name to the appname const
//		2) add the FQDN of your G Suite domain to the testdomain const
//		3) add the ID of the Google sheet to the testsheetid const
//			note: the part after https://docs.google.com/spreadsheets/d/ but before /edit
//		4) set the folowing environment variables to suit your specific needs:
//			export TESTAPP="gsuitemdmtest"
//			export TESTDOMAIN="yourdomain.com"
//			export TESTSHEETID="1bnfhj459dbhs95ngkbnvbnlsjvpas82bhh5d_9W8fjs"
//			export GOOGLE_APPLICATION_CREDENTIALS="/path/to/credentials_yourdomain.com.json"
//		5) go get -u github.com/rickt/gsuitemdm
//		6) go build gsuitemdmtest.go
//		7) ./gsuitemdmtest

var (
	appname     string = os.Getenv("TESTAPP")
	testdomain  string = os.Getenv("TESTDOMAIN")
	testsheetid string = os.Getenv("TESTSHEETID")
	configfile  string = "gsuitemdmtest_conf.json"
)

func main() {
	// get a context
	ctx := context.Background()

	// get a G Suite MDM Service
	gs, err := gsuitemdm.New(ctx, configfile)
	if err != nil {
		log.Fatal("Couldn't get a gsuitemdm service")
	}

	// setup logging
	lg := gs.Log.Logger(appname)

	// get Admin SDK data
	err = gs.GetAdminSDKDevices(testdomain)
	if err != nil {
		lg.Log(logging.Entry{Payload: err, Severity: logging.Error})
	}
	lg.Log(logging.Entry{Payload: "G Suite Admin SDK for domain '" + testdomain + " reports " + strconv.Itoa(len(gs.SDKData.Mobiledevices)) + " mobile devices", Severity: logging.Notice})

	// get sheet data
	err = gs.GetSheetData()
	if err != nil {
		lg.Log(logging.Entry{Payload: err, Severity: logging.Error})
	}
	lg.Log(logging.Entry{Payload: "Google Sheet reports " + strconv.Itoa(len(gs.SheetData)-1) + " rows", Severity: logging.Notice})

	// get datastore data
	err = gs.GetDatastoreDevices()
	if err != nil {
		lg.Log(logging.Entry{Payload: err, Severity: logging.Error})
	}
	lg.Log(logging.Entry{Payload: "Google Datastore reports " + strconv.Itoa(len(gs.DatastoreData)) + " mobile devices", Severity: logging.Notice})

	// merge some data
	md := gs.MergeDatastoreAndSheetData()

	// update the sheet
	err = gs.UpdateSheet(md)
	if err != nil {
		lg.Log(logging.Entry{Payload: "Error updating Sheet: " + fmt.Sprintf("%s", err), Severity: logging.Error})
	}

	// get admin API data again, update datastore
	err = gs.GetAdminSDKDevices(testdomain)
	if err != nil {
		lg.Log(logging.Entry{Payload: err, Severity: logging.Error})
	}
	count, err := gs.UpdateAllDatastoreDevices(testdomain)
	if err != nil {
		lg.Log(logging.Entry{Payload: err, Severity: logging.Error})
	} else {
		lg.Log(logging.Entry{Payload: "Updated Datastore with " + strconv.Itoa(count) + " mobile devices", Severity: logging.Notice})
	}

	// flush the logs
	err = gs.Log.Close()
	if err != nil {
		lg.Log(logging.Entry{Payload: err, Severity: logging.Error})
	}

	return
}
