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
	"time"
)

// Convert a Datastore mobile device object to an Admin SDK mobile device object
func (mdms *GSuiteMDMService) ConvertDatastoreDevicetoSDK(device *DatastoreMobileDevice) *admin.MobileDevice {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

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
func (mdms *GSuiteMDMService) GetDatastoreDevices() error {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

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

// Search for a matching device in Google Datastore using a specific Admin SDK mobile device object
func (mdms *GSuiteMDMService) SearchDatastoreForDevice(device *admin.MobileDevice) (*DatastoreMobileDevice, error) {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

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
		}
	}

	// Return
	return nil, errors.New(fmt.Sprintf("Could not find device: %s", err))
}

// Update a specific device in Google Cloud Datastore
func (mdms *GSuiteMDMService) UpdateDeviceInDatastore(device *admin.MobileDevice) error {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

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
