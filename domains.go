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
func (ms *GSuiteMDMService) buildFullDomainList() []string {
	if ms.GlobalDebug {
		defer TimeTrack(time.Now())
	}

	var domains []string

	// Range through the slice of configured domains
	for _, d := range ms.Domains {
		domains = append(domains, d.DomainName)
	}

	return domains
}

// Get a CustomerID for a given domain
func (ms *GSuiteMDMService) getDomainCustomerID(domain string) string {
	if ms.GlobalDebug {
		defer TimeTrack(time.Now())
	}

	// Range through the slice of configured domains and look for the specified domain
	for _, d := range ms.Domains {
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
	if ms.GlobalDebug {
		defer TimeTrack(time.Now())
	}

	var ok = false

	// Iterate through the slice of configured domains and look for the specified domain
	for _, d := range ms.Domains {
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
	for _, domain := range ms.Domains {
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
	if ms.GlobalDebug {
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
