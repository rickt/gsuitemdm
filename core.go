package gsuitemdm

//
// GSuiteMDM core funcs
//

import (
	"context"
)

// Create a new G Suite MDM Service
func New(ctx context.Context, file string) *GSuiteMDMService {
	// Load in main configuration file and get a config struct
	cf := loadConfigFile(file)

	return &GSuiteMDMService{
		c:   cf,
		ctx: ctx,
	}
}

// EOF
