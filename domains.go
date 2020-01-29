package gsuitemdm

//
// GSuiteMDM domain-specific funcs
//

import (
	"errors"
	"fmt"
)

// Build a list of all configured domains
func (mdms *GSuiteMDMService) BuildFullDomainList() []string {
	var domains []string

	// Range through the slice of configured domains
	for _, d := range mdms.C.Domains {
		domains = append(domains, d.DomainName)
	}

	return domains
}

// Get a CustomerID for a given domain
func (mdms *GSuiteMDMService) GetDomainCustomerID(domain string) (string, error) {
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

// EOF
