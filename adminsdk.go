package gsuitemdm

//
// GSuiteMDM G Suite Admin SDK-specific funcs
//

import (
	admin "google.golang.org/api/admin/directory/v1"
	"strings"
	"time"
)

// Convert an Admin SDK mobile device object to a Datastore mobile device object
func (mdms *GSuiteMDMService) ConvertSDKDeviceToDatastore(device *admin.MobileDevice) (*DatastoreMobileDevice, error) {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

	var d, y DatastoreMobileDevice
	var err error
	var x = new(DatastoreMobileDevice)

	// Convert data received from the Admin SDK
	d.CompromisedStatus = device.DeviceCompromisedStatus
	d.DeveloperMode = device.DeveloperOptionsStatus
	d.Domain = getEmailDomain(device.Email[0])
	d.Email = device.Email[0]
	d.IMEI = strings.Replace(device.Imei, " ", "", -1)
	d.Model = device.Model
	d.Name = device.Name[0]
	d.OS = device.Os
	d.OSBuild = device.BuildNumber
	d.SN = strings.Replace(device.SerialNumber, " ", "", -1)
	d.Status = device.Status
	d.SyncFirst = device.FirstSync
	d.SyncLast = device.LastSync
	d.Type = device.Type
	d.UnknownSources = device.UnknownSourcesStatus
	d.USBADB = device.AdbStatus
	d.WifiMac = device.WifiMacAddress

	// If Datastore has existing "local data" (colour, notes, phone number, RAM) for this device, we need to merge
	// it with the data received from the Admin SDK
	x, err = mdms.SearchDatastoreForDevice(device)
	if err != nil {
		return &d, err
	}

	d.Color = x.Color
	d.Notes = x.Notes
	d.PhoneNumber = x.PhoneNumber
	d.RAM = x.RAM

	// However, if the Google Sheet also has exinsting local data for this device, we need to merge it as well.
	y = mdms.SearchSheetForDevice(device)

	d.Color = y.Color
	d.Notes = y.Notes
	d.PhoneNumber = y.PhoneNumber
	d.RAM = y.RAM

	return &d, nil
}

// Get the list of devices for a G Suite domain from the Admin SDK
func (mdms *GSuiteMDMService) GetAdminSDKDevices(domain string) error {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

	var as *admin.Service
	var cid string
	var err error

	// Iterate through main config struct until we find the specific domain
	for _, d := range mdms.C.Domains {
		switch {
		case d.DomainName == domain:
			// Domain found!
			cid, err = mdms.GetDomainCustomerID(domain)
			if err != nil {
				return err
			}

			// Authenticate with this domain
			as, err = mdms.AuthenticateWithDomain(cid, domain, mdms.C.SearchScope)
			if err != nil {
				return err
			}

			// Pull down the list of devices for this G Suite domain.
			// Refer to https://godoc.org/google.golang.org/api/admin/directory/v1#MobileDevices
			mdms.SDKData, err = as.Mobiledevices.List(d.CustomerID).OrderBy(mdms.C.APIQueryOrderBy).Do()
			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

// EOF
