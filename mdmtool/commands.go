package main

//
// MDMTool commands
//

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

//
// APPROVE
//

// Add the "approve" command
func addApproveCommand(mdmtool *kingpin.Application) {
	c := &ApproveCommand{}
	approve := mdmtool.Command("approve", "Approve a mobile device").Action(c.run)
	approve.Flag("domain", "The G Suite domain to which the mobile device belongs to (required)").Required().Short('d').StringVar(&c.Domain)
	approve.Flag("imei", "Approve using a mobile device IMEI number").Short('i').StringVar(&c.IMEI)
	approve.Flag("sn", "Approve using a mobile device serial number").Short('s').StringVar(&c.SN)
}

// Setup the "approve" command
func (ac *ApproveCommand) run(c *kingpin.ParseContext) error {
	fmt.Printf("approve goes here\n")
	return nil
}

//
// BLOCK
//

// Add the "block" command
func addBlockCommand(mdmtool *kingpin.Application) {
	c := &BlockCommand{}
	block := mdmtool.Command("block", "Block a mobile device").Action(c.run)
	block.Flag("domain", "The G Suite domain to which the mobile device belongs to (required)").Required().Short('d').StringVar(&c.Domain)
	block.Flag("imei", "Block using a mobile device IMEI number").Short('i').StringVar(&c.IMEI)
	block.Flag("sn", "Block using a mobile device serial number").Short('s').StringVar(&c.SN)
}

// Setup the "block" command
func (ac *BlockCommand) run(c *kingpin.ParseContext) error {
	fmt.Printf("block goes here\n")
	return nil
}

//
// DELETE
//

// Add the "delete" command
func addDeleteCommand(mdmtool *kingpin.Application) {
	c := &DeleteCommand{}
	del := mdmtool.Command("delete", "Delete a mobile device").Action(c.run)
	del.Flag("domain", "The G Suite domain to which the mobile device belongs to (required)").Required().Short('d').StringVar(&c.Domain)
	del.Flag("imei", "Delete using a mobile device IMEI number").Short('i').StringVar(&c.IMEI)
	del.Flag("sn", "Delete using a mobile device serial number").Short('s').StringVar(&c.SN)
}

// Setup the "delete" command
func (ac *DeleteCommand) run(c *kingpin.ParseContext) error {
	fmt.Printf("delete goes here\n")
	return nil
}

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

//
// LIST DOMAINS
//

// Add the "listdomains" command
func addListDomainsCommand(mdmtool *kingpin.Application) {
	c := &ListDomainsCommand{}
	ld := mdmtool.Command("listdomains", "List all configured G Suite domains").Action(c.run)
	ld.Flag("verbose", "Enable verbose mode").Short('v').BoolVar(&c.Verbose)

}

// Setup the "listdomains" command
func (ld *ListDomainsCommand) run(c *kingpin.ParseContext) error {
	fmt.Printf("list domains goes here\n")
	return nil
}

//
// WIPE
//

// Add the "wipe" command
func addWipeCommand(mdmtool *kingpin.Application) {
	c := &WipeCommand{}
	wipe := mdmtool.Command("wipe", "Wipe a mobile device").Action(c.run)
	wipe.Flag("domain", "The G Suite domain to which the mobile device belongs to (required)").Required().Short('d').StringVar(&c.Domain)
	wipe.Flag("imei", "Wipe using a mobile device IMEI number").Short('i').StringVar(&c.IMEI)
	wipe.Flag("sn", "Wipe using a mobile device serial number").Short('s').StringVar(&c.SN)
}

// Setup the "wipe" command
func (ac *WipeCommand) run(c *kingpin.ParseContext) error {
	fmt.Printf("wipe goes here\n")
	return nil
}

// EOF
