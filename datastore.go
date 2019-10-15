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

// Read all mobile device data from Google Cloud Datastore
func (mdms *GSuiteMDMService) GetDatastoreDevices() ([]DatastoreMobileDevice, error) {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

	var devices []DatastoreMobileDevice
	var err error

	// Create a Datastore client
	client, err := datastore.NewClient(mdms.Ctx, mdms.C.ProjectID)
	if err != nil {
		return devices, errors.New(fmt.Sprintf("Error creating Datastore client: %s", err))
	}

	// Build the query & get the list of devices
	_, err = client.GetAll(mdms.Ctx, datastore.NewQuery("MobileDevice").
		Order("Name"),
		&devices)
	if err != nil {
		return devices, errors.New(fmt.Sprintf("Error querying Datastore: %s", err))
	}

	// Return
	return devices, nil

}

// Search for a matching device in Google Datastore using a specific Admin SDK mobile device object
func (mdms *GSuiteMDMService) SearchDatastoreForDevice(device *admin.MobileDevice) (*DatastoreMobileDevice, error) {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

	var d = new(DatastoreMobileDevice)
	var dsd []DatastoreMobileDevice
	var err error = nil

	// Get device data from Datastore
	dsd, err = mdms.GetDatastoreDevices()
	if err != nil {
		return d, errors.New(fmt.Sprintf("Error querying Datastore: %s", err))
	}

	// Normalise the IMEI we're looking for
	nimei := strings.Replace(device.Imei, " ", "", -1)

	// Range through the slice of devices from Datastore, and when found, return it
	for k := range dsd {
		if nimei == strings.Replace(dsd[k].IMEI, " ", "", -1) {
			// Found!
			d = &dsd[k]
			break
		}
	}

	// Return
	return d, nil
}

// EOF
