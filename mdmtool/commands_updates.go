package main

//
// MDMTool update commands (updatedb, updatesheet)
//
//

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"net/http"
)

//
// UPDATE DATASTORE
//
func addUpdateDatastoreCommand(mdmtool *kingpin.Application) {
	c := &UpdateDatastoreCommand{}
	ud := mdmtool.Command("updatedb", "Update the DB").Action(c.run)
	ud.Flag("verbose", "Enable verbose mode").Short('v').BoolVar(&c.Verbose)

}

// Setup the "updatedb" command
func (ud *UpdateDatastoreCommand) run(c *kingpin.ParseContext) error {
	var rb gsuitemdm.UpdateRequest

	// Setup the request body
	rb.Key = m.Config.APIKey

	// Marshal the JSON
	js, err := json.Marshal(rb)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updating Datastore... ")

	// Build the http request
	req, err := http.NewRequest("POST", m.Config.UpdateDatastoreURL, bytes.NewBuffer(js))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create an http client
	client := &http.Client{}

	// Send the request and get a nice response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Printf(" done.\n")

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	return nil
}

//
// UPDATE SHEET
//
func addUpdateSheetCommand(mdmtool *kingpin.Application) {
	c := &UpdateSheetCommand{}
	us := mdmtool.Command("updatesheet", "Update the Google Sheet").Action(c.run)
	us.Flag("verbose", "Enable verbose mode").Short('v').BoolVar(&c.Verbose)

}

// Setup the "updatesheet" command
func (us *UpdateSheetCommand) run(c *kingpin.ParseContext) error {
	var rb gsuitemdm.UpdateRequest

	// Setup the request body
	rb.Key = m.Config.APIKey

	// Marshal the JSON
	js, err := json.Marshal(rb)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updating Google Sheet... ")

	// Build the http request
	req, err := http.NewRequest("POST", m.Config.UpdateSheetURL, bytes.NewBuffer(js))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create an http client
	client := &http.Client{}

	// Send the request and get a nice response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Printf(" done.\n")

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	return nil
}

// EOF
