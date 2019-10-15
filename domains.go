package gsuitemdm

//
// GSuiteMDM domain-specific funcs
//

import (
	"errors"
	"fmt"
	"time"
)

// Build a list of all configured domains
func (mdms *GSuiteMDMService) BuildFullDomainList() []string {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

	var domains []string

	// Range through the slice of configured domains
	for _, d := range mdms.C.Domains {
		domains = append(domains, d.DomainName)
	}

	return domains
}

// Get a CustomerID for a given domain
func (mdms *GSuiteMDMService) GetDomainCustomerID(domain string) (string, error) {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

	// Range through the slice of configured domains and look for the specified domain
	for _, d := range mdms.C.Domains {
		switch d.DomainName {
		case domain:
			// Domain found, return the domain's CustomerID
			return d.CustomerID, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Could not find CustomerID for domain %s", domain))
}

// Check to see if a domain is configured
func (mdms *GSuiteMDMService) IsDomainConfigured(domain string) bool {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

	var ok = false

	// Iterate through the slice of configured domains and look for the specified domain
	for _, d := range mdms.C.Domains {
		if domain == d.DomainName {
			// Domain found!
			ok = true
			break
		}
	}

	return ok
}

// List all configured domains
func (mdms *GSuiteMDMService) ListAllDomains(verbose bool) {
	// Range through the slice of configured domains and print out some nice info
	for _, domain := range mdms.C.Domains {
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
func (mdms *GSuiteMDMService) VerifySpecifiedDomain(domain string) ([]string, error) {
	if mdms.C.Debug {
		defer TimeTrack(time.Now())
	}

	var domains []string

	// Check to see if the specified domain is configured within mdmtool
	if mdms.IsDomainConfigured(domain) == false {
		// Domain is not valid
		t := "ERROR: domain '" + domain + "' is not a valid or configured domain"
		return nil, errors.New(t)
	}

	// Domain is valid
	domains = append(domains, domain)

	return domains, nil
}

// EOF
