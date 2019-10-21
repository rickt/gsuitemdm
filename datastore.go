package gsuitemdm

//
// GSuiteMDM Google Datastore-specific funcs
//

import (
	"cloud.google.com/go/datastore"
	"errors"
	"fmt"
	admin "google.golang.org/api/admin/directory/v1"
	"log"
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
	// TODO: change MobileDevice below
	_, err = dc.GetAll(mdms.Ctx, datastore.NewQuery("MobileDevice").
		Order("Name"),
		&mdms.DatastoreData)

	if err != nil {
		return errors.New(fmt.Sprintf("Error querying Datastore: %s", err))
	}

	// Return
	return nil
}

// Search for a matching device in Google Datastore using a specific Admin SDK mobile device object
func (mdms *GSuiteMDMService) SearchDatastoreForDevice(device *admin.MobileDevice) (*DatastoreMobileDevice, error) {
	var d = new(DatastoreMobileDevice)
	var err error

	// Normalise the IMEI we're looking for
	nimei := strings.Replace(device.Imei, " ", "", -1)

	log.Printf("SearchDatastoreForDevice(): looking for device=%s\n", nimei)

	// Range through the slice of devices from Datastore, and when found, return it
	for k := range mdms.DatastoreData {
		if nimei == strings.Replace(mdms.DatastoreData[k].IMEI, " ", "", -1) {
			// Found!
			log.Printf("SearchDatastoreForDevice(): device found, device=%v\n", device)
			d = &mdms.DatastoreData[k]
			return d, nil
		} else {
			// Not found, lets create it
			log.Printf("SearchDatastoreForDevice(): device NOT found, device=%v\n", device)
			d, err = mdms.ConvertSDKDeviceToDatastore(device)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("SearchDatastoreForDevice(): 1Could not find device: %s, device=%v", err, device))
			}
			return d, nil
		}
	}

	// Return
	return nil, errors.New(fmt.Sprintf("SearchDatastoreForDevice(): 2Could not find device: %s, device=%v", err, device))
}

// Update a device in Google Cloud Datastore
func (mdms *GSuiteMDMService) UpdateDatastoreDevice(device *admin.MobileDevice) error {
	var ed = new(DatastoreMobileDevice)
	var nd = new(DatastoreMobileDevice)
	var dc *datastore.Client
	var err error
	var key *datastore.Key

	// Create a Datastore client
	dc, err = datastore.NewClient(mdms.Ctx, mdms.C.ProjectID)
	if err != nil {
		return err
	}

	// We were passed an Admin SDK mobile device object. We need to convert it to a
	// new Datastore mobile device object
	nd, err = mdms.ConvertSDKDeviceToDatastore(device)
	if err != nil {
		return err
	}

	// Get the existing Datastore entry for this device
	key = datastore.NameKey(mdms.C.DSNamekey, nd.SN, nil)
	err = dc.Get(mdms.Ctx, key, ed)
	if err != nil {
		return err
	}

	// If existing data exists for this device in Datastore, preserve it
	if ed.PhoneNumber != "" {
		nd.PhoneNumber = strings.Replace(ed.PhoneNumber, " ", "", -1)
	}
	if ed.Color != "" {
		nd.Color = ed.Color
	}
	if ed.RAM != "" {
		nd.RAM = ed.RAM
	}
	if ed.Notes != "" {
		nd.Notes = ed.Notes
	}

	// If existing data exists for this device in the Google Sheet, preserve it
	for _, shv := range mdms.SheetData {
		if (strings.Replace(nd.IMEI, " ", "", -1) == strings.Replace(shv.IMEI, " ", "", -1)) ||
			(strings.Replace(nd.SN, " ", "", -1) == strings.Replace(shv.SN, " ", "", -1)) {
			nd.Color = shv.Color
			nd.RAM = shv.RAM
			nd.Notes = shv.Notes
			nd.PhoneNumber = strings.Replace(shv.PhoneNumber, " ", "", -1)
		}
	}

	// We're finished, save the device in Datastore
	_, err = dc.Put(mdms.Ctx, key, nd)

	if err != nil {
		return err
	}

	return nil
}

// EOF
