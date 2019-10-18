package gsuitemdm

//
// GSuiteMDM Google Datastore-specific funcs
//

import (
	"cloud.google.com/go/datastore"
	"errors"
	"fmt"
	admin "google.golang.org/api/admin/directory/v1"
	"strings"
)

// Convert a Datastore mobile device object to an Admin SDK mobile device object
func (mdms *GSuiteMDMService) ConvertDatastoreDevicetoSDK(device *DatastoreMobileDevice) *admin.MobileDevice {
	var d admin.MobileDevice

	// Convert
	d.DeviceCompromisedStatus = device.CompromisedStatus
	d.DeveloperOptionsStatus = device.DeveloperMode
	d.Email[0] = device.Email
	d.Imei = strings.Replace(device.IMEI, " ", "", -1)
	d.Model = device.Model
	d.Os = device.OS
	d.BuildNumber = device.OSBuild
	d.SerialNumber = strings.Replace(device.SN, " ", "", -1)
	d.Status = device.Status
	d.FirstSync = device.SyncFirst
	d.LastSync = device.SyncLast
	d.Type = device.Type
	d.UnknownSourcesStatus = device.UnknownSources
	d.AdbStatus = device.USBADB
	d.WifiMacAddress = device.WifiMac

	return &d
}

// Read all mobile device data from Google Cloud Datastore
func (mdms *GSuiteMDMService) GetDatastoreData() error {
	var dc *datastore.Client
	var err error

	// Create a Datastore client
	dc, err = datastore.NewClient(mdms.Ctx, mdms.C.ProjectID)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating Datastore client: %s", err))
	}

	// Build the query & get the list of devices
	_, err = dc.GetAll(mdms.Ctx, datastore.NewQuery("MobileDevice").
		Order("Name"),
		&mdms.DatastoreData)

	if err != nil {
		return errors.New(fmt.Sprintf("Error querying Datastore: %s", err))
	}

	// Return
	return nil
}

// Merge Datastore mobile device data & Google Sheet mobile device data
func (mdms *GSuiteMDMService) MergeDatastoreAndSheetData() []DatastoreMobileDevice {
	var mergeddata []DatastoreMobileDevice

	// Merge the data
	for _, dsv := range mdms.DatastoreData {
		// Create a temporary mobile device using data from datastore
		var d DatastoreMobileDevice

		// Merge the data
		d.CompromisedStatus = dsv.CompromisedStatus
		d.Domain = dsv.Domain
		d.DeveloperMode = dsv.DeveloperMode
		d.Email = dsv.Email
		d.IMEI = dsv.IMEI
		d.Model = dsv.Model
		d.Name = dsv.Name
		d.OS = dsv.OS
		d.OSBuild = dsv.OSBuild
		d.SN = dsv.SN
		d.Status = dsv.Status
		d.SyncFirst = dsv.SyncFirst
		d.SyncLast = dsv.SyncLast
		d.Type = dsv.Type
		d.UnknownSources = dsv.UnknownSources
		d.USBADB = dsv.USBADB
		d.WifiMac = dsv.WifiMac

		// Add the local-to-sheet data for this specific mobile device (if it exists)
		for _, shv := range mdms.SheetData {
			if (strings.Replace(d.IMEI, " ", "", -1) == strings.Replace(shv.IMEI, " ", "", -1)) ||
				(strings.Replace(d.SN, " ", "", -1) == strings.Replace(shv.SN, " ", "", -1)) {
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

// Search for a matching device in Google Datastore using a specific Admin SDK mobile device object
func (mdms *GSuiteMDMService) SearchDatastoreForDevice(device *admin.MobileDevice) (*DatastoreMobileDevice, error) {
	var d = new(DatastoreMobileDevice)
	var err error

	// Normalise the IMEI we're looking for
	nimei := strings.Replace(device.Imei, " ", "", -1)

	// Range through the slice of devices from Datastore, and when found, return it
	for k := range mdms.DatastoreData {
		if nimei == strings.Replace(mdms.DatastoreData[k].IMEI, " ", "", -1) {
			// Found!
			d = &mdms.DatastoreData[k]
			return d, nil
		} else {
			// Not found, lets create it
			d, err = mdms.ConvertSDKDeviceToDatastore(device)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Could not find device: %s", err))
			}
			return d, nil
		}
	}

	// Return
	return nil, errors.New(fmt.Sprintf("Could not find device: %s", err))
}

// Update all Datastore devices for a given domain with device data from the Admin SDK
func (mdms *GSuiteMDMService) UpdateAllDatastoreData() (int, error) {
	var count int
	var d = new(DatastoreMobileDevice)
	var dc *datastore.Client
	var err error

	// Create a Datastore client
	dc, err = datastore.NewClient(mdms.Ctx, mdms.C.ProjectID)
	if err != nil {
		return 0, err
	}

	// Iterate through the domain's devices
	for _, device := range mdms.SDKData.Mobiledevices {

		fmt.Printf("UpdateAllDatastoreData(): device = %v\n", device)

		// Convert our *admin.MobileDevice to an *hmsMobileDevice
		d, err = mdms.ConvertSDKDeviceToDatastore(device)
		if err != nil {
			return 0, err
		}

		fmt.Printf("UpdateAllDatastoreData(): d = %v\n", d)

		// Setup the Datastore key
		key := datastore.NameKey(mdms.C.DSNamekey, d.SN, nil)

		// Save the entity in Datastore
		_, err = dc.Put(mdms.Ctx, key, d)

		if err != nil {
			return 0, err
		}

		// Increment the counter
		count++
	}

	return count, err
}

// Update a specific device in Google Cloud Datastore
func (mdms *GSuiteMDMService) UpdateDatastoreDevice(device *admin.MobileDevice) error {
	var d = new(DatastoreMobileDevice)
	var dc *datastore.Client
	var err error

	// Create a Datastore client
	dc, err = datastore.NewClient(mdms.Ctx, mdms.C.ProjectID)
	if err != nil {
		return err
	}

	// We were passed an Admin SDK mobile device object. We need to convert it to
	// a Datastore mobile device object
	d, err = mdms.ConvertSDKDeviceToDatastore(device)
	if err != nil {
		return err
	}

	// Setup the datastore key
	key := datastore.NameKey(mdms.C.DSNamekey, d.SN, nil)

	// Save the device in Datastore
	_, err = dc.Put(mdms.Ctx, key, d)

	if err != nil {
		return err
	}

	return nil
}

// EOF
