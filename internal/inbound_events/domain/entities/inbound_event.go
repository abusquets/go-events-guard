package entities

import (
	"eventsguard/internal/utils/entities"
	"time"
)

type Retry struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at" validate:"required"`
	Result    string    `json:"result" bson:"result"`
}

type InboundEvent struct {
	ID        entities.ID `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Type      string      `json:"type" bson:"code" validate:"required"`
	ClientID  entities.ID `json:"client_id" bson:"client_id" validate:"required"`
	Version   string      `json:"version" bson:"version" validate:"required"`
	Payload   string      `json:"payload" bson:"payload" validate:"required"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at" validate:"required"`
	SendAt    *time.Time  `json:"send_at" bson:"created_at" validate:"required"`
	Retries   []Retry     `json:"retries" bson:"retries"`
}

func NewInboundEvent(eventType string, clientID entities.ID, version string, payload string, sendAt *time.Time) (*InboundEvent, error) {
	event := &InboundEvent{
		ID:        entities.NewID(),
		Type:      eventType,
		ClientID:  clientID,
		Version:   version,
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
		SendAt:    sendAt,
		Retries:   []Retry{},
	}

	return event, nil
}
