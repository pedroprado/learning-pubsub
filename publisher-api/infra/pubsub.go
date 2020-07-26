package infra

import (
	"context"
	"fmt"
	"os"
	"test/pubsub/models"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
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

func init() {
	ctx = context.Background()
	client = CreateClient(ctx, projectID)
}

//CreatePushSubscription initializes a subscription for pushing messages
func CreatePushSubscription(context *gin.Context) {

	sub, err := generatePushSubscription(client, pushTopicID, pushSubscriptionID)

	if err != nil {
		apiErr := models.ApiError{
			Msg: err.Error(),
		}
		fmt.Println(err)
		context.JSON(500, gin.H{"status": "Internal Error", "error": apiErr})
		return
	}

	fmt.Printf("Subscription created: %v\n", sub)
	context.JSON(200, gin.H{"status": "OK"})
}

func generatePushSubscription(client *pubsub.Client, topicID string, subsID string) (*pubsub.Subscription, error) {

	topic, topicErr := client.CreateTopic(ctx, topicID)
	if topicErr != nil {
		return nil, topicErr
	}

	sub, subErr := client.CreateSubscription(ctx, subsID, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 10 * time.Second,
		PushConfig: pubsub.PushConfig{
			Endpoint: "http://localhost:4000/messages",
		},
	})

	if subErr != nil {
		return nil, subErr
	}

	return sub, nil
}

//CreatePullSubscription initializes a subscription for pulling messages
func CreatePullSubscription() {
	sub, err := generatePullSubscription(client, pullTopicID, pullSubscriptionID)

	if err != nil {
		apiErr := models.ApiError{
			Msg: err.Error(),
		}
		fmt.Printf("%+v \n", apiErr)
		return
	}
	fmt.Printf("Pull Subscription created: %v\n", sub)
}

func generatePullSubscription(client *pubsub.Client, topicID string, subsID string) (*pubsub.Subscription, error) {

	topic, topicErr := client.CreateTopic(ctx, topicID)
	if topicErr != nil {
		return nil, topicErr
	}

	sub, subErr := client.CreateSubscription(ctx, subsID, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})

	if subErr != nil {
		return nil, subErr
	}

	return sub, nil
}

//GetTopics retrieves the topic of a given client
func GetTopics(client *pubsub.Client) ([]*pubsub.Topic, error) {
	var topics []*pubsub.Topic
	it := client.Topics(ctx)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Next: %v", err)
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

//GetPullSubscriptions returns all the subscriptions of a project
func GetPullSubscriptions(client *pubsub.Client, topicID string) ([]*pubsub.Subscription, error) {

	var subs []*pubsub.Subscription

	it := client.Topic(topicID).Subscriptions(ctx)
	for {
		sub, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Next: %v", err)
		}
		subs = append(subs, sub)
	}
	return subs, nil

}

// func GetEnv() {
// 	fmt.Println("Project id: ", projectID)
// 	fmt.Println("PullTopic id: ", pullTopicID)

// 	fmt.Println("PushTopic id: ", pushTopicID)

// }
