package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	projectID              = os.Getenv("PUBSUB_PROJECT_ID")
	pullTopicID            = os.Getenv("PUBSUB_CNC_TOPIC_ID")
	pullSubscriptionID     = os.Getenv("PUBSUB_CNC_SUBSCRIPTION_ID")
	pushTopicID            = os.Getenv("PUBSUB_NOTIFY_TOPIC_ID")
	pushSubscriptionID     = os.Getenv("PUBSUB_NOTIFY_SUBSCRIPTION_ID")
	deadLetterTopicID      = os.Getenv("PUBSUB_DL_TOPIC_ID")
	deadLetterSubscription = os.Getenv("PUBSUB_DL_SUBSCRIPTION_ID")
	ctx                    context.Context
	client                 *pubsub.Client
)

type cncCommand struct {
	EntityID   uuid.UUID `json:"entityId"`
	EngineName string    `json:"engineName"`
	Time       time.Time `json:"time"`
	NoCache    bool      `json:"noCache"`
}

type message struct {
	Data      []byte `json:"data,omitempty"`
	MessageID string `json:"messageId"`
	// Attributes struct{} `json:"attributes"`
}

type pubSubMessage struct {
	Message      message `json:"message"`
	Subscription string  `json:"subscription"`
}

func main() {
	ctx = context.Background()
	client = createClient()

	go pullMessages()

	router := gin.Default()
	router.GET("/dead-letter/pull", deadLetterPull)
	router.Run(":4000")

}

func receivePushMessages(context *gin.Context) {
	// message := pubSubMessage{}
	// bs, _ := ioutil.ReadAll(context.Request.Body)
	// fmt.Printf("%s\n", string(bs))
	// err := json.Unmarshal(bs, &message)
	// if err != nil {
	// 	fmt.Printf("%+v \n", err)
	// }
	// fmt.Printf("%+v \n", message)

	// content := cncCommand{}
	// err2 := json.Unmarshal(message.Message.Data, &content)
	// if err2 != nil {
	// 	fmt.Printf("%+v \n", err2)
	// }
	// fmt.Printf("%+v \n", content)
	// context.JSON(http.StatusOK, gin.H{"status": 200, "data": "OK"})
}

func deadLetterPull(ginCtx *gin.Context) {
	sub := client.Subscription(deadLetterSubscription)
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = -1
	sub.ReceiveSettings.MaxExtension = 0
	sub.ReceiveSettings.MaxExtensionPeriod = 0

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Got message: %q\n", string(msg.Data))
		ginCtx.JSON(http.StatusOK, gin.H{"status": 200, "data": string(msg.Data)})
		msg.Ack()
	})

	if err != nil {
		ginCtx.JSON(500, gin.H{"msg": "Error receiving from deadletter queue", "error": err.Error()})
	}
}

//------------------------------------------------------------------
func pullMessages() {
	channel := make(chan *pubsub.Message)

	go pulling(channel)

	go func() {
		for msg := range channel {
			cnc := cncCommand{}
			bs := msg.Data
			err := json.Unmarshal(bs, &cnc)
			if err != nil {
				logrus.Errorf("Error on unmarshal message: %+v", errors.WithStack(err))
			} else {
				logrus.Printf("Got message with id: %v\n", msg.ID)
				logrus.Printf("%+v\n", cnc)

				err := process(&cnc)
				if err != nil {
					logrus.WithField("message_id=", msg.ID).
						Errorf("Error on processing message: %+v", errors.WithStack(err))
				} else {
					// msg.Ack()
					fmt.Printf("Processed without ack\n")
				}
			}

		}
	}()
}

func pulling(channel chan *pubsub.Message) {
	sub := client.Subscription(pullSubscriptionID)
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = -1
	sub.ReceiveSettings.MaxExtension = 0
	sub.ReceiveSettings.MaxExtensionPeriod = 0

	for {
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			channel <- msg
			time.Sleep(3 * time.Second)
		})
		if err != nil {
			logrus.Errorf("Error to receive message: %+v\n", errors.WithStack(err))
			time.Sleep(10 * time.Second)
		}
	}

}

func process(cnc *cncCommand) error {
	return nil
}

//-------------------------------------------------------
func convertMsg(bs []byte) {
	message := pubSubMessage{}

	err := json.Unmarshal(bs, &message)
	if err != nil {
		fmt.Printf("%+v \n", err)
	}

	content := cncCommand{}
	err2 := json.Unmarshal(message.Message.Data, &content)
	if err2 != nil {
		fmt.Printf("%+v \n", err2)
	}
}

func createClient() *pubsub.Client {

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed do connect client: %v", err)
	}

	fmt.Printf("Client connected without error!\n")
	return client
}
