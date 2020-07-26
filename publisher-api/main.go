package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"test/pubsub/infra"
	"test/pubsub/models"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	projectID          = os.Getenv("PUBSUB_PROJECT_ID")
	pullTopicID        = os.Getenv("PUBSUB_CNC_TOPIC_ID")
	pullSubscriptionID = os.Getenv("PUBSUB_CNC_SUBSCRIPTION_ID")
	pushTopicID        = os.Getenv("PUBSUB_NOTIFY_TOPIC_ID")
	pushSubscriptionID = os.Getenv("PUBSUB_NOTIFY_SUBSCRIPTION_ID")
	deadLetterTopicID  = os.Getenv("PUBSUB_DL_TOPIC_ID")
	ctx                context.Context
	client             *pubsub.Client
)

func main() {

	router := gin.Default()
	ctx = context.Background()
	client = infra.CreateClient(ctx, projectID)
	fmt.Println()
	router.GET("/message/publish", pushMessage)
	router.GET("/dead-letter/create", createDeadLetterTopic)
	router.Run(":5000")
}

func pushMessage(context *gin.Context) {

	cncCommand := models.CncCommand{
		EntityID:   uuid.New(),
		EngineName: "profile",
		Time:       time.Now(),
		NoCache:    true,
	}

	message, err := json.Marshal(&cncCommand)

	if err != nil {
		apiErr := models.ApiError{
			Msg: err.Error(),
		}
		fmt.Println(err)
		context.JSON(500, gin.H{"status": "Internal Error", "error": apiErr})
		return
	}

	publishMessage(client, pullTopicID, message)

	context.JSON(200, gin.H{"status": "OK", "message": "published message"})

}

func createDeadLetterTopic(context *gin.Context) {

	fullyQualifiedDeadLetterTopic := "project/" + projectID + "/topics/" + deadLetterTopicID

	topic := client.Topic(pullTopicID)

	sub, err := client.CreateSubscription(ctx, pullSubscriptionID, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 10 * time.Second,
		DeadLetterPolicy: &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     fullyQualifiedDeadLetterTopic,
			MaxDeliveryAttempts: 5,
		},
	})

	fmt.Println(err)
	if err != nil {
		context.JSON(500, gin.H{"msg": "Error creating deadletter queue", "error": err.Error()})
	} else {
		context.JSON(http.StatusCreated, gin.H{"status": 201, "data": sub})
	}

}

//----------------------------------------------------

func publishMessage(client *pubsub.Client, topicID string, msg []byte) {
	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})

	id, err := result.Get(ctx)
	if err != nil {
		fmt.Println("Error publishing")
	}
	fmt.Println("Published message to the Topic, with messageID: ", id)
}
