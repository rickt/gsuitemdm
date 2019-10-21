package gsuitemdm

//
// GSuiteMDM Google Sheet-specific funcs
//

import (
	"errors"
	"fmt"
	"github.com/Iwark/spreadsheet"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Read all mobile device data from the Google Sheet
func (mdms *GSuiteMDMService) GetSheetData() error {
	var err error

	// Get an authenticated http client
	client, err := mdms.HttpClient(mdms.C.SheetCreds)
	if err != nil {
		return err
	}

	// Get a google sheets service
	gss := spreadsheet.NewServiceWithClient(client)

	// Fetch the Google sheet
	sheet, err := gss.FetchSpreadsheet(mdms.C.SheetID)
	if err != nil {
		return err
	}

	// Select the first worksheet
	ws, err := sheet.SheetByIndex(0)

	// Row count is number of rows in the worksheet minus 1 (because of header row)
	numrows := len(ws.Rows) - 1

	// Range through the Sheet's rows, skipping the first row because it's a header row
	var row = 0
	for kr := range ws.Rows {
		if row == numrows {
			break
		}

		// Build a temporary device
		var d DatastoreMobileDevice

		// Get data from the sheet and populate
		d.PhoneNumber = (ws.Rows[kr+1][1].Value)
		d.Color = (ws.Rows[kr+1][2].Value)
		d.RAM = (ws.Rows[kr+1][3].Value)
		d.IMEI = (ws.Rows[kr+1][8].Value)
		d.SN = (ws.Rows[kr+1][9].Value)
		d.Notes = (ws.Rows[kr+1][17].Value)

		// Append this device to devices
		mdms.SheetData = append(mdms.SheetData, d)

		// Increment the row count
		row++
	}

	return nil
}

// Create an authenticated http(s) client, used to read/write the Google Sheet
func (mdms *GSuiteMDMService) HttpClient(creds string) (*http.Client, error) {
	// Read in the JSON credentials file for the domain/user we will write the Google Sheet as
	data, err := ioutil.ReadFile(creds)
	if err != nil {
		return nil, err
	}

	// Get a nice juicy JWT config struct using that credentials file
	conf, err := google.JWTConfigFromJSON(data, mdms.C.SheetScope)
	if err != nil {
		return nil, err
	}

	// Since we are using a service account's JSON credentials to write, we need to specify
	// an actual G Suite user (required by Google)
	conf.Subject = mdms.C.SheetWho

	// Return the authenticated http client
	return conf.Client(oauth2.NoContext), nil
}

// Merge Datastore and Sheet data
func (mdms *GSuiteMDMService) MergeDatastoreAndSheetData() []DatastoreMobileDevice {
	var mergeddata []DatastoreMobileDevice

	// Range through the Datastore data
	for _, dsv := range mdms.DatastoreData {
		// Create a temporary mobile device using data from Datastore
		var d DatastoreMobileDevice

		// Merge
		d.CompromisedStatus = dsv.CompromisedStatus
		d.Domain = dsv.Domain
		d.DeveloperMode = dsv.DeveloperMode
		d.Email = dsv.Email
		d.IMEI = strings.Replace(dsv.IMEI, " ", "", -1)
		d.Model = dsv.Model
		d.Name = dsv.Name
		d.OS = dsv.OS
		d.OSBuild = dsv.OSBuild
		d.SN = strings.Replace(dsv.SN, " ", "", -1)
		d.Status = dsv.Status
		d.SyncFirst = dsv.SyncFirst
		d.SyncLast = dsv.SyncLast
		d.Type = dsv.Type
		d.UnknownSources = dsv.UnknownSources
		d.USBADB = dsv.USBADB
		d.WifiMac = dsv.WifiMac

		// Add the local-to-sheet data for this specific mobile device (if it exists)
		for _, shv := range mdms.SheetData {
			if (strings.Replace(d.IMEI, " ", "", -1) == strings.Replace(shv.IMEI, " ", "", -1)) || (strings.Replace(d.SN, " ", "", -1) == strings.Replace(shv.SN, " ", "", -1)) {
				d.Color = shv.Color
				d.RAM = shv.RAM
				d.Notes = shv.Notes
				d.PhoneNumber = shv.PhoneNumber
			}
		}

		// Append this mobile device to the slice of mobile devices
		mergeddata = append(mergeddata, d)
	}

	return mergeddata
}

// Search the Google Sheet for a specific device
func (mdms *GSuiteMDMService) SearchSheetForDevice(device *admin.MobileDevice) (DatastoreMobileDevice, error) {
	var d DatastoreMobileDevice

	// Add the local-to-Sheet data for this specific mobile device (if it exists)
	for _, shv := range mdms.SheetData {
		if (strings.Replace(device.Imei, " ", "", -1) == strings.Replace(shv.IMEI, " ", "", -1)) ||
			(strings.Replace(device.SerialNumber, " ", "", -1) == strings.Replace(shv.SN, " ", "", -1)) {
			// Device found!
			return shv, nil
		}
	}

	return d, errors.New(fmt.Sprintf("Could not find device"))
}

// Update the Google Sheet
func (mdms *GSuiteMDMService) UpdateSheet(mergeddata []DatastoreMobileDevice) error {
	// Get an authenticated http client
	client, err := mdms.HttpClient(mdms.C.SheetCreds)
	if err != nil {
		return err
	}

	// Get a Google Sheets service
	gss := spreadsheet.NewServiceWithClient(client)

	// Fetch the Google sheet
	sheet, err := gss.FetchSpreadsheet(mdms.C.SheetID)
	if err != nil {
		return err
	}

	// Select the first worksheet
	ws, err := sheet.SheetByIndex(0)

	// Note that we start at row 2 because row0 == "Last updated" line and row1 == Header
	var row = 2

	// Set time zone to be as configured
	loc, err := time.LoadLocation(mdms.C.TimeZone)
	if err != nil {
		return err
	}
	// Update the Last Updated timestamp in the sheet
	ws.Update(0, 1, time.Now().In(loc).Format(time.RFC1123))

	// Range through the canonical device data
	for _, upd := range mergeddata {
		// Update each column, per row
		ws.Update(row, 0, upd.Domain)
		ws.Update(row, 1, upd.PhoneNumber)
		ws.Update(row, 2, upd.Color)
		ws.Update(row, 3, upd.RAM)
		ws.Update(row, 4, upd.Name)
		ws.Update(row, 5, upd.Status)
		ws.Update(row, 6, upd.Email)
		ws.Update(row, 7, upd.Model)
		ws.Update(row, 8, upd.IMEI)
		ws.Update(row, 9, upd.SN)
		ws.Update(row, 10, upd.SyncLast)
		ws.Update(row, 11, upd.OS)
		ws.Update(row, 12, upd.Type)
		ws.Update(row, 13, upd.WifiMac)
		ws.Update(row, 14, upd.CompromisedStatus)
		ws.Update(row, 15, strconv.FormatBool(upd.DeveloperMode))
		ws.Update(row, 16, strconv.FormatBool(upd.UnknownSources))
		ws.Update(row, 17, strconv.FormatBool(upd.USBADB))
		ws.Update(row, 18, upd.Notes)

		// Incremement the row count
		row++
	}

	// Save all changes to the Sheet
	err = ws.Synchronize()
	if err != nil {
		return err
	}

	// Return
	return nil
}

// EOF
