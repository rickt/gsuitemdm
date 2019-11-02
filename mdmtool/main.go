package main

//
// MDMTool main code
//

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

var (
	appname    string = "mdmtool"
	configfile string = "mdmtool_conf.json"
	m                 = new(MDMTool)
)

func main() {
	var err error

	// Load MDMTool config
	m.Config, err = loadMDMToolConfig(configfile)
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
