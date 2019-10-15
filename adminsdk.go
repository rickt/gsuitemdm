package gsuitemdm

//
// GSuiteMDM G Suite Admin SDK-specific funcs
//

import (
	admin "google.golang.org/api/admin/directory/v1"
	"time"
)

// Get the list of devices for a G Suite domain from the Admin API
func (mdms *GSuiteMDMService) GetAdminAPIDeviceData(domain string) (*admin.MobileDevices, error) {
	if mdms.C.GlobalDebug {
		defer TimeTrack(time.Now())
	}

	var ad *admin.MobileDevices
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
				return ad, err
			}

			// Authenticate with this domain
			as, err = mdms.AuthenticateWithDomain(cid, domain, mdms.C.SearchScope)
			if err != nil {
				return ad, err
			}

			// Pull down the list of devices for this G Suite domain
			// The List method we're calling returns a MobileDevices.
			// Refer to https://godoc.org/google.golang.org/api/admin/directory/v1#MobileDevices
			devices, err := as.Mobiledevices.List(d.CustomerID).OrderBy(mdms.C.APIQueryOrderBy).Do()
			if err != nil {
				return ad, err
			}

			return devices, nil
		}
	}

	return nil, nil
}

// EOF
