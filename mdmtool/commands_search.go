package main

//
// MDMTool commands
//
//

import (
	"fmt"
	// "github.com/rickt/gsuitemdm"
	"gopkg.in/alecthomas/kingpin.v2"
	// "io/ioutil"
	// "log"
	// "net/http"
)

//
// DIRECTORY
//

// Add the "directory" command
func addDirectoryCommand(mdmtool *kingpin.Application) {
	c := &DirectoryCommand{}
	dir := mdmtool.Command("dir", "Search the mobile phone directory").Action(c.run)
	dir.Flag("name", "Search using name").Short('n').StringVar(&c.Name)
	dir.Flag("email", "Search using email").Short('e').StringVar(&c.Email)
}

// Setup the "directory" command
func (ac *DirectoryCommand) run(c *kingpin.ParseContext) error {
	fmt.Printf("directory goes here\n")
	return nil
}

// EOF
