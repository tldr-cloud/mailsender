package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	errors "errors"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func accessSecretVersion(secretChannel chan string, errorChannel chan error, name string) {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		errorChannel <- errors.New(fmt.Sprint("failed to create secretmanager client: %v", err))
		return
	}

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		errorChannel <- errors.New(fmt.Sprint("failed to access secret version: %v", err))
		return
	}

	secretChannel <- string(result.Payload.Data)
}

//func main() {
//	secretChannel := make(chan string)
//	errorChannel := make(chan error)
//
//	secretName := "projects/346496881273/secrets/postmarkapp-server-API-token/versions/latest"
//	go accessSecretVersion(secretChannel, errorChannel, secretName)
//	select {
//		case secret := <-secretChannel:
//			fmt.Print(secret)
//		case err := <-errorChannel:
//			log.Println(err.Error())
//	}
//}
