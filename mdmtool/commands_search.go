package main

//
// MDMTool search commands (dir, search)
//
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
	"sort"
)

//
// DIRECTORY
//

// Add the "directory" command
func addDirectoryCommand(mdmtool *kingpin.Application) {
	c := &DirectoryCommand{}
	dir := mdmtool.Command("dir", "Search the mobile phone directory").Action(c.run)
	dir.Flag("name", "Search for a phone number using name").Short('n').StringVar(&c.Name)
	dir.Flag("email", "Search for a phone number using email").Short('e').StringVar(&c.Email)
}

// Setup the "directory" command
func (dr *DirectoryCommand) run(c *kingpin.ParseContext) error {
	// Check runtime options
	if (dr.Email == "" && dr.Name == "") || (dr.Email != "" && dr.Name != "") {
		return errors.New("with \"dir\" command you must specify one of --email or --name")
	}

	// Runtime options are good, lets setup the request body
	var rb gsuitemdm.SearchRequest

	// What kind of search are we doing?
	switch {

	// Email
	case dr.Email != "":
		rb.QType = "email"
		rb.Q = dr.Email
		break

	// Name
	case dr.Name != "":
		rb.QType = "name"
		rb.Q = dr.Name
		break
	}

	// Setup the rest of the SEARCH request
	rb.Key = m.Config.APIKey

	// Marshal the JSON
	js, err := json.Marshal(rb)
	if err != nil {
		log.Fatal(err)
	}

	// Build the http request
	req, err := http.NewRequest("POST", m.Config.DirectoryURL, bytes.NewBuffer(js))
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON
	var dirdata []gsuitemdm.DirectoryData
	err = json.Unmarshal(body, &dirdata)

	// If this was a bad request, or no results returned, exit
	if len(dirdata) < 1 {
		// Was this a bad request?
		if resp.Status == "400 Bad Request" {
			fmt.Printf("%s\n", body)
		}
		if resp.Status == "204 No Content" {
			// Or was this a good response but just with no data?
			fmt.Printf("Search returned 0 results.\n")
		}
		return nil
	}

	// Okay, we have good data, sort it
	sort.Sort(gsuitemdm.AllDirectoryData{dirdata})

	// Print a nice header line
	printDirectoryHeaderLine()

	// Range through the directory entries and pretty-print then
	for k := range dirdata {
		printDirectoryData(dirdata[k])
	}

	// Print a nice footer line
	printDirectoryLine()

	fmt.Printf("Search returned %d results.\n", len(dirdata))

	return nil
}

//
// SEARCH
//

// Add the "search" command
func addSearchCommand(mdmtool *kingpin.Application) {
	c := &SearchCommand{}
	search := mdmtool.Command("search", "Search for mobile devices").Action(c.run)
	search.Flag("all", "Show all mobile devices").Short('a').BoolVar(&c.All)
	search.Flag("domain", "Restrict search to a specific G Suite domain (optional)").Short('d').StringVar(&c.Domain)
	search.Flag("email", "Search for a device using email address").Short('e').StringVar(&c.Email)
	search.Flag("imei", "Search for a device using IMEI").Short('i').StringVar(&c.IMEI)
	search.Flag("name", "Search for a device using staff name").Short('n').StringVar(&c.Name)
	search.Flag("notes", "Search for a device using notes").Short('o').StringVar(&c.Notes)
	search.Flag("phone", "Search for a device using phone number").Short('p').StringVar(&c.Phone)
	search.Flag("sn", "Search for a device using serial number").Short('s').StringVar(&c.SN)
	search.Flag("status", "Search for a device using MDM device status").Short('t').StringVar(&c.Status)
	search.Flag("verbose", "Enable verbose mode").Short('v').BoolVar(&c.Verbose)
}

// Setup the "search" command
func (sc *SearchCommand) run(c *kingpin.ParseContext) error {
	// Check runtime options
	if sc.All != true && sc.Email == "" && sc.IMEI == "" && sc.Name == "" && sc.Notes == "" && sc.Phone == "" && sc.SN == "" && sc.Status == "" {
		return errors.New("with \"search\" command you must specify one of --all, --email, --imei, --name, --phone, --sn or --status")
	}

	// Check runtime options: cannot use other search operators when using --all
	if sc.All == true && (sc.Email != "" || sc.IMEI != "" || sc.Name != "" || sc.Notes != "" || sc.Phone != "" || sc.SN != "" || sc.Status != "") {
		return errors.New("with \"search --all\" you cannot also specify --email, --imei, --name, --phone, --sn or --status")
	}

	// Runtime options are good, lets setup the request body
	var rb gsuitemdm.SearchRequest

	// What kind of search are we doing?
	switch {

	// All
	case sc.All == true:
		rb.QType = "all"
		break

	// Email
	case sc.Email != "":
		rb.QType = "email"
		rb.Q = sc.Email
		break

	// IMEI
	case sc.IMEI != "":
		rb.QType = "imei"
		rb.Q = sc.IMEI
		break

	// Name
	case sc.Name != "":
		rb.QType = "name"
		rb.Q = sc.Name
		break

	// Notes
	case sc.Notes != "":
		rb.QType = "notes"
		rb.Q = sc.Notes
		break

	// Phone
	case sc.Phone != "":
		rb.QType = "phone"
		rb.Q = sc.Phone
		break

	// Serial Number
	case sc.SN != "":
		rb.QType = "sn"
		rb.Q = sc.SN
		break

	// Status
	case sc.Status != "":
		rb.QType = "status"
		rb.Q = sc.Status
		break
	}

	// Setup the rest of the SEARCH request
	rb.Domain = sc.Domain
	rb.Key = m.Config.APIKey

	// Marshal the JSON
	js, err := json.Marshal(rb)
	if err != nil {
		log.Fatal(err)
	}

	// Build the http request
	req, err := http.NewRequest("POST", m.Config.SearchDatastoreURL, bytes.NewBuffer(js))
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON
	var reply []gsuitemdm.DatastoreMobileDevice
	err = json.Unmarshal(body, &reply)

	// If this was a bad request, or no results returned, exit
	if len(reply) < 1 {
		// Was this a bad request?
		if resp.Status == "400 Bad Request" {
			fmt.Printf("%s\n", body)
		}
		if resp.Status == "204 No Content" {
			// Or was this a good response but just with no data?
			fmt.Printf("Search returned 0 results.\n")
		}
		return nil
	}

	// Okay, we have good data, sort it
	sort.Sort(gsuitemdm.DatastoreMobileDevices{reply})

	// Only print header line if verbose mode was NOT requested
	if sc.Verbose != true {
		printHeaderLine()
	}

	// Range through the returned data and pretty-print it
	for k := range reply {
		printDeviceData(reply[k], sc.Verbose)
	}

	// Only print final line if verbose mode was NOT requested
	if sc.Verbose != true {
		printLine()
	}

	fmt.Printf("Search returned %d results.\n", len(reply))

	return nil
}

// EOF
