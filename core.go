package gsuitemdm

//
// GSuiteMDM core funcs
//

import (
	"context"
)

// Create a new G Suite MDM Service
func New(ctx context.Context, file string) (*GSuiteMDMService, error) {
	var cf GSuiteMDMConfig
	var err error

	// Load in main configuration file and get a config struct
	cf, err = loadConfig(file)
	if err != nil {
		return nil, err
	}

	// Create a new G Suite MDM service and populate it
	return &GSuiteMDMService{
		C:   cf,
		Ctx: ctx}, nil
}

// EOF
