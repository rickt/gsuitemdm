package gsuitemdm

//
// GSuiteMDM Google Datastore-specific funcs
//

import (
	"cloud.google.com/go/datastore"
	"errors"
	"fmt"
	"time"
)

// Read all mobile device data from Google Cloud Datastore
func (mdms *GSuiteMDMService) GetDatastoreDevices() ([]DatastoreMobileDevice, error) {
	if mdms.C.GlobalDebug {
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

// EOF
