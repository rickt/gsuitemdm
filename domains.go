package gsuitemdm

//
// GSuiteMDM domain-specific funcs
//

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"io/ioutil"
	"log"
	"time"
)

// Authenticate with a domain, get an admin.Service
func (ms *GSuiteMDMService) authenticateWithDomain(customerid string, domain string, scope string) *admin.Service {
	if ms.Config.GlobalDebug {
		defer TimeTrack(time.Now())
	}

	// Range through slice of configured domains until we find the domain we're looking for
	for _, d := range ms.Config.Domains {
		switch {
		// Domain found!
		case d.DomainName == domain:
			// Read in this domain's service account JSON credentials file
			creds, err := ioutil.ReadFile(d.CredentialsFile)
			checkError(err)

			// create JWT config using the credentials file
			jwt, err := google.JWTConfigFromJSON(creds, scope)
			checkError(err)

			// Specify which admin user the API calls should "run as"
			jwt.Subject = d.AdminUser

			// Make the API client using our JWT config
			as, err := admin.New(jwt.Client(context.Background()))
			checkError(err)

			// Return the admin.Service
			return as
		}
	}

	log.Fatalf("Error: could not authenticate with domain %s\n", domain)
	return nil
}

// Build a list of all configured domains
func (ms *GSuiteMDMService) buildFullDomainList() []string {
	if ms.Config.GlobalDebug {
		defer TimeTrack(time.Now())
	}

	var domains []string

	// Range through the slice of configured domains
	for _, d := range ms.Config.Domains {
		domains = append(domains, d.DomainName)
	}

	return domains
}

// Get a CustomerID for a given domain
func (ms *GSuiteMDMService) getDomainCustomerID(domain string) string {
	if ms.Config.GlobalDebug {
		defer TimeTrack(time.Now())
	}

	// Range through the slice of configured domains and look for the specified domain
	for _, d := range ms.Config.Domains {
		switch d.DomainName {
		case domain:
			// Domain found, return the domain's CustomerID
			return d.CustomerID
		}
	}

	return ""
}

// Check to see if a domain is configured
func (ms *GSuiteMDMService) isDomainConfigured(domain string) bool {
	if ms.Config.GlobalDebug {
		defer TimeTrack(time.Now())
	}

	var ok = false

	// Iterate through the slice of configured domains and look for the specified domain
	for _, d := range ms.Config.Domains {
		if domain == d.DomainName {
			// Domain found!
			ok = true
			break
		}
	}

	return ok
}

// List all configured domains
func (ms *GSuiteMDMService) listAllDomains(verbose bool) {
	// Range through the slice of configured domains and print out some nice info
	for _, domain := range ms.Config.Domains {
		if verbose == true {
			fmt.Printf("%s:\n", domain.DomainName)
			fmt.Printf("	customerid: %s\n", domain.CustomerID)
			fmt.Printf("	credentialsfile: %s\n", domain.CredentialsFile)
			fmt.Printf("	adminuser: %s\n", domain.AdminUser)
		} else {
			fmt.Printf("	%s\n", domain.DomainName)
		}
	}
}

// Verify specified domain
func (ms *GSuiteMDMService) verifySpecifiedDomain(domain string) ([]string, error) {
	if ms.Config.GlobalDebug {
		defer TimeTrack(time.Now())
	}

	var domains []string

	// Check to see if the specified domain is configured within mdmtool
	if ms.isDomainConfigured(domain) == false {
		// Domain is not valid
		t := "ERROR: domain '" + domain + "' is not a valid or configured domain"
		return nil, errors.New(t)
	}

	// Domain is valid
	domains = append(domains, domain)

	return domains, nil
}

// EOF
