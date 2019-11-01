package main

//
// MDMTool commands
//

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"net/http"
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
	// Check runtime options
	if (ac.IMEI == "" && ac.SN == "") || (ac.IMEI != "" && ac.SN != "") {
		return errors.New("with \"approve\" command you must specify either --imei or --sn")
	}

	// Runtime options are good, lets setup the request body
	var approval bool
	var rb gsuitemdm.ActionRequest

	// How are we identifying the device to be approved?
	switch {

	// IMEI
	case ac.IMEI != "":
		rb.IMEI = ac.IMEI
		approval = checkUserConfirmation(fmt.Sprintf("WARNING: Are you sure you want to APPROVE device IMEI=%s in domain %s?", ac.IMEI, ac.Domain))
		break

	// Serial Number
	case ac.SN != "":
		rb.SN = ac.SN
		approval = checkUserConfirmation(fmt.Sprintf("WARNING: Are you sure you want to APPROVE device SN=%s in domain %s?", ac.SN, ac.Domain))
		break
	}

	// Check if approval was given
	if approval == false {
		fmt.Printf("Approval not given, exiting\n")
		return nil

	}

	// Approval has been given, lets setup the rest of the APPROVE request
	rb.Action = "approve"
	rb.Confirm = true
	rb.Domain = ac.Domain
	rb.Key = m.Config.APIKey

	// Marshal the JSON
	js, err := json.Marshal(rb)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Build the http request
	req, err := http.NewRequest("POST", m.Config.ApproveDeviceURL, bytes.NewBuffer(js))
	req.Header.Set("Content-Type", "application/json")

	// Create an http client
	client := &http.Client{}

	// Send the request and get a nice response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

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
func (bc *BlockCommand) run(c *kingpin.ParseContext) error {
	// Check runtime options
	if (bc.IMEI == "" && bc.SN == "") || (bc.IMEI != "" && bc.SN != "") {
		return errors.New("with \"block\" command you must specify either --imei or --sn")
	}

	// Runtime options are good, lets setup the request body
	var approval bool
	var rb gsuitemdm.ActionRequest

	// How are we identifying the device to be blocked?
	switch {

	// IMEI
	case bc.IMEI != "":
		rb.IMEI = bc.IMEI
		approval = checkUserConfirmation(fmt.Sprintf("WARNING: Are you sure you want to BLOCK device IMEI=%s in domain %s?", bc.IMEI, bc.Domain))
		break

	// Serial Number
	case bc.SN != "":
		rb.SN = bc.SN
		approval = checkUserConfirmation(fmt.Sprintf("WARNING: Are you sure you want to BLOCK device SN=%s in domain %s?", bc.SN, bc.Domain))
		break
	}

	// Check if approval was given
	if approval == false {
		fmt.Printf("Approval not given, exiting\n")
		return nil

	}

	// Approval has been given, lets setup the rest of the BLOCK request
	rb.Action = "block"
	rb.Confirm = true
	rb.Domain = bc.Domain
	rb.Key = m.Config.APIKey

	// Marshal the JSON
	js, err := json.Marshal(rb)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Build the http request
	req, err := http.NewRequest("POST", m.Config.BlockDeviceURL, bytes.NewBuffer(js))
	req.Header.Set("Content-Type", "application/json")

	// Create an http client
	client := &http.Client{}

	// Send the request and get a nice response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

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
// UPDATE DATASTORE
//
func addUpdateDatastoreCommand(mdmtool *kingpin.Application) {
	c := &UpdateDatastoreCommand{}
	ld := mdmtool.Command("updatedb", "Update the DB").Action(c.run)
	ld.Flag("verbose", "Enable verbose mode").Short('v').BoolVar(&c.Verbose)

}

// Setup the "updatedb" command
func (ld *UpdateDatastoreCommand) run(c *kingpin.ParseContext) error {
	fmt.Printf("updatedb goes here\n")
	return nil
}

//
// UPDATE SHEET
//
func addUpdateSheetCommand(mdmtool *kingpin.Application) {
	c := &UpdateSheetCommand{}
	ld := mdmtool.Command("updatesheet", "Update the Google Sheet").Action(c.run)
	ld.Flag("verbose", "Enable verbose mode").Short('v').BoolVar(&c.Verbose)

}

// Setup the "updatesheet" command
func (ld *UpdateSheetCommand) run(c *kingpin.ParseContext) error {
	fmt.Printf("updatesheet goes here\n")
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
