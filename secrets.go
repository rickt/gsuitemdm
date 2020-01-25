package gsuitemdm

//
// GSuiteMDM GCP Secret Manager related code
//

import (
	"context"
	"errors"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
)

// Get a secret from Secret Manager
func getSecret(ctx context.Context, sid string) (string, error) {
	// Create a Secret Manager client
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", errors.New("Error creating Secret Manage client: " + err.Error())
	}

	// Build the Secret Manager request
	smreq := &secretmanagerpb.AccessSecretVersionRequest{
		Name: sid,
	}

	// Call the Secret Manager API and get the requested secret using its ID
	smres, err := client.AccessSecretVersion(ctx, smreq)
	if err != nil {
		return "", errors.New("Error retrieving secret: " + err.Error())
	}

	// Return the specified secret
	return string(smres.Payload.Data), nil
}

// EOF
