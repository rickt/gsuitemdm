package main

//
// MDMTool helper funcs
//

import (
	"bufio"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/rickt/gsuitemdm"
	"github.com/ttacon/libphonenumber"
	"log"
	"os"
	"strings"
	"time"
)

// Helper function to ask for user confirmation
func checkUserConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// Load MDMTool configuration
func loadMDMToolConfig() (MDMToolConfig, error) {
	c := MDMToolConfig{
		APIKey:             apikey,
		ApproveDeviceURL:   approvedeviceurl,
		BlockDeviceURL:     blockdeviceurl,
		DeleteDeviceURL:    deletedeviceurl,
		DirectoryURL:       directoryurl,
		SearchDatastoreURL: searchdatastoreurl,
		ShowDomainsURL:     showdomainsurl,
		UpdateDatastoreURL: updatedatastoreurl,
		UpdateSheetURL:     updatesheeturl,
		WipeDeviceURL:      wipedeviceurl,
	}

	return c, nil
}

// Print out mobile device data (Datastore edition)
func printDeviceData(device gsuitemdm.DatastoreMobileDevice, verbose bool) {

	// convert last sync strings to time.Time so we can humanize them
	lts, err := time.Parse(time.RFC3339, device.SyncLast)
	if err != nil {
		log.Fatal(err)
	}

	// Make telephone numbers pretty again
	var num *libphonenumber.PhoneNumber
	var fnum string
	if device.PhoneNumber != "" {
		num, err = libphonenumber.Parse(device.PhoneNumber, "US")
		fnum = libphonenumber.Format(num, libphonenumber.NATIONAL)
	}

	// Print detail only if --verbose was specified
	switch verbose {
	case false:
		fmt.Printf("%21.21s | %-16.16s | %-14.14s | %-16.16s | %-15.15s | %-13.13s | %-18.18s | %-20.20s\n", device.Domain, device.Model, fnum, device.SN, device.IMEI, device.Status, humanize.Time(lts), device.Name)

	case true:
		// convert last sync strings to time.Time so we can humanize them
		fts, err := time.Parse(time.RFC3339, device.SyncFirst)
		if err != nil {
			log.Fatal(err)
		}

		// Print the device details
		fmt.Printf("\n")
		fmt.Printf("         Phone Number: %s\n", fnum)
		fmt.Printf("        Device Domain: %s\n", device.Domain)
		fmt.Printf(" Device Serial Number: %s\n", device.SN)
		fmt.Printf("          Device IMEI: %s\n", device.IMEI)
		fmt.Printf("          Device Type: %s\n", device.Type)
		fmt.Printf("        Device Status: %s\n", device.Status)
		fmt.Printf("      Device Wifi Mac: %s\n", device.WifiMac)
		fmt.Printf("    Device Model & OS: %s (%s, build %s)\n", device.Model, device.OS, device.OSBuild)
		fmt.Printf("        Color/Storage: %s /  %s\n", device.Color, device.RAM)
		fmt.Printf("           Owner Name: %s\n", device.Name)
		fmt.Printf("          Owner Email: %s\n", device.Email)
		fmt.Printf("           First Sync: %s\n", humanize.Time(fts))
		fmt.Printf("            Last Sync: %s\n", humanize.Time(lts))
		fmt.Printf("   Compromised Status: %s\n", device.CompromisedStatus)
		fmt.Printf("    Encryption Status: %s\n", device.EncryptionStatus)
		fmt.Printf("           OS Options: Developer mode (%v), Allow Unknown Sources (%v), USB Debugging (%v)\n", device.DeveloperMode, device.UnknownSources, device.USBADB)
		fmt.Printf("            --- Notes: ---\n%s\n            --- Notes: ---\n", device.Notes)

	}
	return
}

// Print out directory data
func printDirectoryData(person gsuitemdm.DirectoryData) {
	fmt.Printf("%-21.21s | %s | %s\n", person.Name, person.PhoneNumber, person.Email)

	// fmt.Printf("%21.21s | %-16.16s | %-14.14s | %-16.16s | %-15.15s | %-13.13s | %-18.18s | %-20.20s\n", device.Domain, device.Model, fnum, device.SN, device.IMEI, device.Status, humanize.Time(lts), device.Name)
}

// Print out the directory header line
func printDirectoryHeaderLine() {
	// print the first line of dashes
	printDirectoryLine()
	// print header line
	fmt.Printf("Name                  | Phone Number   | Email \n")
	// print a line of dashes under the header line
	printDirectoryLine()
}

// Print out a correctly formatted line for the phone directory
func printDirectoryLine() {
	// print a line
	fmt.Printf("----------------------+----------------+------------------------------------------\n")
}

// Print out the header
func printHeaderLine() {
	// print the first line of dashes
	printLine()
	// print header line
	fmt.Printf("Domain                | Model            | Phone Number   | Serial #         | IMEI            | Status        | Last Sync          | Owner\n")
	// print a line of dashes under the header line
	printLine()
}

// Print a correctly formatted line
func printLine() {
	// print a line
	fmt.Printf("----------------------+------------------+----------------+------------------+-----------------+---------------+--------------------+---------------\n")
}

// EOF
