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
//		* downloads all mobile device data for a configured G Suite domain's MDM-managed devices using the G Suite Admin SDK API
//		* downloads all mobile device data from a configured tracking Google sheet (if data already exists in the sheet)
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
//		1) add your own app name to the 'appname' var
//		2) add the FQDN of your G Suite domain to the 'testdomain' var
//		3) edit the 'configfile', making sure to change the following vars to suit:
//			'globaldebug': if you want debug messages (bool)
//			'projectid': set to the name of the GCP project you want to run gsuitemdm inside
//			'sheetcreds': set to the path of the JSON credentials file of the user with
//				appropriate permissions to write the test Google sheet
//			'sheetid': set to the ID of the test Google sheet. Note this is the part
//				after https://docs.google.com/spreadsheets/d/ but before /edit
//			'sheetwho': set to the email address of the G Suite user who has permissions
//				to write the test Google sheet
//			'domains': setup this JSON array as per your G Suite domain setup
//		4) set the folowing environment variables to suit your specific needs:
//			export TESTAPP="gsuitemdmtest"
//			export TESTDOMAIN="yourdomain.com"
//			export GOOGLE_APPLICATION_CREDENTIALS="/path/to/credentials_yourdomain.com.json"
//		5) go get -u github.com/rickt/gsuitemdm
//		6) go build
//		7) ./gsuitemdmtest
//
// Implementation Notes:
//    * Google Sheet:
//       An assumption is made that the Google Sheet expects data to start on row 3.
//       Row 1 will be autofilled with "Last Automatic Update: <timestamp>"
//       Row 2 is a header row, with columns A through R as (in CSV) format:
//          DOMAIN,WIRELESS #,COLOR,SIZE,OWNER NAME,STATUS,OWNER EMAIL,MODEL,IMEI/ESN Hex,Serial #,\
//          OS,TYPE,WIFI MAC,COMPROMISED?,DEV MODE?,UNKNOWN SOURCES?,USB DEBUG?,NOTES

var (
	appname    string = os.Getenv("TESTAPP")
	testdomain string = os.Getenv("TESTDOMAIN")
	configfile string = "gsuitemdmtest_conf.json"
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
