package gsuitemdm

//
// GSuiteMDM core funcs
//

import (
	"context"
)

// Create a new G Suite MDM Service
func (mdms *GSuiteMDMService) New(ctx context.Context) *GSuiteMDMService {
	return &GSuiteMDMService{
		ctx: ctx,
	}
}

// EOF
