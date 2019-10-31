package main

import (
	"context"
	"github.com/rickt/gsuitemdm"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

// MDMTool main code

var (
	appname    string = "mdmtool"
	configfile string = "mdmtool_conf.json"
	urlsfile   string = "mdmtool_urls.json"
)

func main() {
	var err error
	var m = new(MDMTool)

	// Get a context
	ctx := context.Background()

	// Get a G Suite MDM Service
	m.GSMDMService, err = gsuitemdm.New(ctx, configfile)
	if err != nil {
		log.Fatal("Couldn't get a gsuitemdm service")
	}

	// Load MDMTool URLs configuration
	m.URLs, err = loadMDMToolURLs(urlsfile)
	if err != nil {
		log.Fatal("Error loading MDMTool URLs")
	}

	// Create an MDMTool
	mdmtool := kingpin.New("mdmtool", "HMS G Suite MDM Tool")

	// Add the commands
	addApproveCommand(mdmtool)     // approve
	addBlockCommand(mdmtool)       // block
	addDeleteCommand(mdmtool)      // delete
	addDirectoryCommand(mdmtool)   // directory
	addListDomainsCommand(mdmtool) // list domains
	addWipeCommand(mdmtool)        // wipe

	// Parse runtime options
	kingpin.MustParse(mdmtool.Parse(os.Args[1:]))

}

// EOF
