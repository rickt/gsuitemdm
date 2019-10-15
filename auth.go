package gsuitemdm

//
// GSuiteMDM authentication-specific funcs
//

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"io/ioutil"
	"net/http"
	"time"
)

// Authenticate with a domain, get an admin.Service
func (mdms *GSuiteMDMService) AuthenticateWithDomain(customerid, domain string) (*admin.Service, error) {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

	// Range through slice of configured domains until we find the domain we're looking for
	for _, d := range mdms.C.Domains {
		switch {
		// Domain found!
		case d.DomainName == domain:
			// Read in this domain's service account JSON credentials file
			creds, err := ioutil.ReadFile(d.CredentialsFile)
			if err != nil {
				return nil, err
			}

			// create JWT config using the credentials file
			jwt, err := google.JWTConfigFromJSON(creds, mdms.C.SearchScope)
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

// Create an authenticated http(s) client, used to read/write the Google Sheet
func (mdms *GSuiteMDMService) HttpClient(creds string) (*http.Client, error) {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

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

// EOF
