package main

import (
	"context"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"log"
	"os"
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
//		1) edit the 'configfile', making sure to change the following vars to suit:
//			'globaldebug': if you want debug messages (bool)
//			'projectid': set to the name of the GCP project you want to run gsuitemdm inside
//			'sheetcreds': set to the path of the JSON credentials file of the user with
//				appropriate permissions to write the test Google sheet
//			'sheetid': set to the ID of the test Google sheet. Note this is the part
//				after https://docs.google.com/spreadsheets/d/ but before /edit
//			'sheetwho': set to the email address of the G Suite user who has permissions
//				to write the test Google sheet
//			'domains': setup this JSON array as per your G Suite domain setup
//		2) set the folowing environment variables to suit your specific needs:
//			export TESTAPP="gsuitemdmtest"
//			export TESTDOMAIN="yourdomain.com"
//			export TESTSHEETID="1bnfhj459dbhs95ngkbnvbnlsjvpas82bhh5d_9W8fjs"
//			export GOOGLE_APPLICATION_CREDENTIALS="/path/to/credentials_yourdomain.com.json"
//		3) go get -u github.com/rickt/gsuitemdm
//		4) go build
//		5) ./gsuitemdmtest

var (
	appname    string = os.Getenv("TESTAPP")
	configfile string = "gsuitemdmtest_conf.json"
	domain     string = os.Getenv("TESTDOMAIN")
)

func main() {
	// Get a context
	ctx := context.Background()

	// Get a G Suite MDM Service
	gs, err := gsuitemdm.New(ctx, configfile)
	if err != nil {
		log.Fatal("Couldn't get a gsuitemdm service")
	}

	// Get Admin SDK data
	err = gs.GetAdminSDKDevices(domain)
	if err != nil {
		fmt.Printf("Error getting mobile device data from G Suite Admin SDK: %v\n", err)
		return
	}
	fmt.Printf("G Suite Admin SDK for domain %s reports %d mobile devices\n", domain, len(gs.SDKData.Mobiledevices))

	// Get Google Sheet data
	err = gs.GetSheetData()
	if err != nil {
		fmt.Printf("Error getting Google Sheet data: %v\n", err)
		return
	}
	fmt.Printf("Google Sheet reports %d rows\n", len(gs.SheetData)-1)

	// Get Datastore data
	err = gs.GetDatastoreData()
	if err != nil {
		fmt.Printf("Error getting Google Datastore data: %v\n", err)
		return
	}
	fmt.Printf("Google Datastore reports %d mobile devices\n", len(gs.DatastoreData))

	// Merge some data
	md := gs.MergeDatastoreAndSheetData()

	// Update the sheet
	err = gs.UpdateSheet(md)
	if err != nil {
		fmt.Printf("Error updating Google Sheet: %v\n", err)
		return
	}

	// Range through the slice of configured domains
	for _, dm := range gs.C.Domains {
		domain = dm.DomainName

		// Get data about this domain's devices from the Admin SDK
		err = gs.GetAdminSDKDevices(domain)
		if err != nil {
			fmt.Printf("Error getting Admin SDK data for domain %s: %s\n", domain, err)
			return
		}

		// Range through this domain's list of devices and update it in Datastore
		for _, device := range gs.SDKData.Mobiledevices {
			err = gs.UpdateDatastoreDevice(device)
			if err != nil {
				fmt.Printf("Error updating a mobile device for domain %s: %s\n", domain, err)
			}
		}
	}

	return
}

// EOF
