package infra

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
)

var (
	serverPort             = os.Getenv("SERVER_PORT")
	projectID              = os.Getenv("PUBSUB_PROJECT_ID")
	pullTopicID            = os.Getenv("PUBSUB_CNC_TOPIC_ID")
	pullSubscriptionID     = os.Getenv("PUBSUB_CNC_SUBSCRIPTION_ID")
	pushTopicID            = os.Getenv("PUBSUB_NOTIFY_TOPIC_ID")
	pushSubscriptionID     = os.Getenv("PUBSUB_NOTIFY_SUBSCRIPTION_ID")
	deadLetterTopicID      = os.Getenv("PUBSUB_DL_TOPIC_ID")
	deadLetterSubscription = os.Getenv("PUBSUB_DL_SUBSCRIPTION_ID")
	listenerInstance       = os.Getenv("LISTENER_INSTANCE")
	ctx                    context.Context
	client                 *pubsub.Client
)
