package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

// MDMTool main code

var (
	appname    string = "mdmtool"
	configfile string = "mdmtool_conf.json"
)

func main() {
	var err error
	var m = new(MDMTool)

	// Load MDMTool config
	m.Config, err = loadMDMToolConfig(configfile)
	if err != nil {
		log.Fatal("Error loading MDMTool URLs")
	}

	// Create an MDMTool
	mdmtool := kingpin.New("mdmtool", "HMS G Suite MDM Tool")

	// Add the commands
	addApproveCommand(mdmtool)         // approve
	addBlockCommand(mdmtool)           // block
	addDeleteCommand(mdmtool)          // delete
	addDirectoryCommand(mdmtool)       // directory
	addListDomainsCommand(mdmtool)     // listdomains
	addUpdateDatastoreCommand(mdmtool) // updatedb
	addUpdateSheetCommand(mdmtool)     // updatesheet
	addWipeCommand(mdmtool)            // wipe

	// Parse runtime options
	kingpin.MustParse(mdmtool.Parse(os.Args[1:]))

}

// EOF
