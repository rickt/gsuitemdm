package main

//
// MDMTool main code
//

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

const (
	apikey             string = "XXXX"
	appname            string = "mdmtool"
	approvedeviceurl   string = "https://us-central1-PROJECTID.cloudfunctions.net/ApproveDevice"
	blockdeviceurl     string = "https://us-central1-PROJECTID.cloudfunctions.net/BlockDevice"
	deletedeviceurl    string = "https://us-central1-PROJECTID.cloudfunctions.net/DeleteDevice"
	directoryurl       string = "https://us-central1-PROJECTID.cloudfunctions.net/Directory"
	searchdatastoreurl string = "https://us-central1-PROJECTID.cloudfunctions.net/SearchDatastore"
	updatedatastoreurl string = "https://us-central1-PROJECTID.cloudfunctions.net/UpdateDatastore"
	updatesheeturl     string = "https://us-central1-PROJECTID.cloudfunctions.net/UpdateSheet"
	wipedeviceurl      string = "https://us-central1-PROJECTID.cloudfunctions.net/WipeDevice"
)

var (
	m = new(MDMTool)
)

func main() {
	var err error

	// Load MDMTool config
	m.Config, err = loadMDMToolConfig()
	if err != nil {
		log.Fatal("Error loading MDMTool configuration")
	}

	// Create an MDMTool app
	mdmtool := kingpin.New(appname, "A G Suite MDM Tool")

	// Add the commands
	addApproveCommand(mdmtool)         // approve
	addBlockCommand(mdmtool)           // block
	addDeleteCommand(mdmtool)          // delete
	addDirectoryCommand(mdmtool)       // directory
	addSearchCommand(mdmtool)          // search
	addUpdateDatastoreCommand(mdmtool) // updatedb
	addUpdateSheetCommand(mdmtool)     // updatesheet
	addWipeCommand(mdmtool)            // wipe

	// Parse runtime options
	kingpin.MustParse(mdmtool.Parse(os.Args[1:]))

}

// EOF
