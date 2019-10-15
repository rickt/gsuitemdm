package gsuitemdm

//
// GSuiteMDM Google Sheet-specific funcs
//

import (
	"github.com/Iwark/spreadsheet"
	admin "google.golang.org/api/admin/directory/v1"
	"strconv"
	"strings"
	"time"
)

// Read all mobile device data from the Google Sheet
func (mdms *GSuiteMDMService) GetSheetData() error {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

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

// Search the Google Sheet for a specific device
func (mdms *GSuiteMDMService) SearchSheetForDevice(device *admin.MobileDevice) DatastoreMobileDevice {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

	var d DatastoreMobileDevice

	// Add the local-to-Sheet data for this specific mobile device (if it exists)
	for _, shv := range mdms.SheetData {
		if (strings.Replace(device.Imei, " ", "", -1) == strings.Replace(shv.IMEI, " ", "", -1)) ||
			(strings.Replace(device.SerialNumber, " ", "", -1) == strings.Replace(shv.SN, " ", "", -1)) {
			// Device found!
			return shv
		}
	}

	return d
}

// Update the Google Sheet
func (mdms *GSuiteMDMService) UpdateSpreadsheet(mergeddata []DatastoreMobileDevice) error {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

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

	// Note that we start at row 2 because row0 == "Last updated" line and row1 == Header
	var row = 2

	// Update the Last Updated timestamp.
	// We want time expressed as being local to Los Angeles
	loc, _ := time.LoadLocation("America/Los_Angeles")
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
		ws.Update(row, 10, upd.OS)
		ws.Update(row, 11, upd.Type)
		ws.Update(row, 12, upd.WifiMac)
		ws.Update(row, 13, upd.CompromisedStatus)
		ws.Update(row, 14, strconv.FormatBool(upd.DeveloperMode))
		ws.Update(row, 15, strconv.FormatBool(upd.UnknownSources))
		ws.Update(row, 16, strconv.FormatBool(upd.USBADB))
		ws.Update(row, 17, upd.Notes)

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
