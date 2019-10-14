package gsuitemdm

//
// GSuiteMDM core funcs
//

import (
	"cloud.google.com/go/logging"
	"context"
)

// Create a new G Suite MDM Service
func New(ctx context.Context, file string) *GSuiteMDMService {
	// Load in main configuration file and get a config struct
	cf := loadConfig(file)

	// Logging (Stackdriver)
	logsd, err := logging.NewClient(ctx, cf.ProjectID)
	checkError(err)

	// Return a new G Suite MDM service
	return &GSuiteMDMService{
		C:   cf,
		Ctx: ctx,
		Log: logsd,
	}
}

// EOF
