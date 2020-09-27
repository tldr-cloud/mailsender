package p

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	errors "errors"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func accessSecretVersion(secretName string) (string, error) {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to create secretmanager client: %v", err))
	}

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to access secret version: %v", err))
	}

	return string(result.Payload.Data), nil
}
