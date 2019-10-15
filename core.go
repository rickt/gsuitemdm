package gsuitemdm

//
// GSuiteMDM core funcs
//

import (
	"cloud.google.com/go/logging"
	"context"
)

// Create a new G Suite MDM Service
func New(ctx context.Context, file string) (*GSuiteMDMService, error) {
	var cf GSuiteMDMConfig
	var err error
	var log *logging.Client

	// Load in main configuration file and get a config struct
	cf, err = loadConfig(file)
	if err != nil {
		return nil, err
	}

	// Setup logging
	log, err = logging.NewClient(ctx, cf.ProjectID)
	if err != nil {
		return nil, err
	}

	// Create a new G Suite MDM service and populate it
	return &GSuiteMDMService{
		C:   cf,
		Ctx: ctx,
		Log: log}, nil
}

// EOF
