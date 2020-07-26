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

type ApiError struct {
	Msg string `json:"message"`
}
