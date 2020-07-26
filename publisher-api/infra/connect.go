package infra

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

//CreateClient creates an instance of a connection to a Pubsub project
func CreateClient(ctx context.Context, projectID string) *pubsub.Client {

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed do connect client: %v", err)
	}

	fmt.Printf("Client connected without error!\n")
	return client
}
