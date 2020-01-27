package gsuitemdm

//
// GSuiteMDM G Suite Admin SDK-specific funcs
//

import (
	"context"
	"errors"
	"fmt"
	"github.com/rickt/gsuitemdm"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
)

// Authenticate with a domain, get an admin.Service
func (mdms *GSuiteMDMService) AuthenticateWithDomain(customerid, domain, scope string) (*admin.Service, error) {
	// Range through slice of configured domains until we find the domain we're looking for
	for _, d := range mdms.C.Domains {
		switch {
		// Domain found!
		case d.DomainName == domain:
			// Get credentials for this domain from Secret Manager. First, get a
			// context and a Secret Manager client
			ctx := context.Background()
			client, err := secretmanager.NewClient(ctx)
			if err != nil {
				return nil, err
			}

			// Build the Secret Manager request
			smreq := &secretmanagerpb.AccessSecretVersionRequest{
				Name: d.SecretID,
			}

			// Call the Secret Manager API and retrieve the ID of the configuration secret
			smres, err := client.AccessSecretVersion(ctx, smreq)
			if err != nil {
				return nil, err
			}
			
			// Retrieve this domain's configuration from Secret Manager
	    config, err := gsuitemdm.GetSecret(ctx, smres.Payload.Data)
	    if err != nil {
        return nil, err
	    }

			// create JWT config using the credentials file
			jwt, err := google.JWTConfigFromJSON(config, scope)
			if err != nil {
				return nil, err
			}

			// Specify which admin user the API calls should "run as"
			jwt.Subject = d.AdminUser

			// Make the API client using our JWT config
			as, err := admin.New(jwt.Client(context.Background()))
			if err != nil {
				return nil, err
			}

			// We've made it all the way through (w00t!), so return the admin.Service
			return as, nil
		}
	}

	// trombone.wav
	return nil, errors.New(fmt.Sprintf("Could not authenticate with domain %s", domain))
}

// Convert an Admin SDK mobile device object to a Datastore mobile device object
func (mdms *GSuiteMDMService) ConvertSDKDeviceToDatastore(device *admin.MobileDevice) (*DatastoreMobileDevice, error) {
	var d DatastoreMobileDevice

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
	d.ResourceId = device.ResourceId
	d.SN = strings.Replace(device.SerialNumber, " ", "", -1)
	d.Status = device.Status
	d.SyncFirst = device.FirstSync
	d.SyncLast = device.LastSync
	d.Type = device.Type
	d.UnknownSources = device.UnknownSourcesStatus
	d.USBADB = device.AdbStatus
	d.WifiMac = device.WifiMacAddress

	return &d, nil
}

// Get the list of devices for a G Suite domain from the Admin SDK
func (mdms *GSuiteMDMService) GetAdminSDKDevices(domain string) error {
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

			// Pull down the list of devices for this G Suite domain
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
