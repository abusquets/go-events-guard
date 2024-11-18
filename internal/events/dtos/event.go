package dtos

import (
	"eventsguard/internal/utils/entities"
	"time"
)

type CreateEventInput struct {
	Type     string      `json:"type" validate:"required"`
	ClientID entities.ID `json:"client_id" validate:"required"`
	Version  string      `json:"version" validate:"required"`
	Payload  string      `json:"payload" validate:"required"`
	SendAt   *time.Time  `json:"send_at,omitempty"`
}
