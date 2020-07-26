package models

import (
	"time"

	"github.com/google/uuid"
)

type CncCommand struct {
	EntityID   uuid.UUID `json:"entityId"`
	EngineName string    `json:"engineName"`
	Time       time.Time `json:"time"`
	NoCache    bool      `json:"noCache"`
}

type Message struct {
	Data      []byte `json:"data,omitempty"`
	MessageID string `json:"messageId"`
	// Attributes struct{} `json:"attributes"`
}

type PubSubMessage struct {
	Message      Message `json:"message"`
	Subscription string  `json:"subscription"`
}
