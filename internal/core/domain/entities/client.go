package entities

import (
	"eventsguard/internal/utils/entities"
)

type Client struct {
	ID       entities.ID `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Code     string      `json:"code" bson:"code" validate:"required"`
	Name     string      `json:"name" bson:"name" validate:"required"`
	IsActive bool        `json:"is_active" bson:"is_active" validate:"required"`
}

func NewClient(code string, name string, isActive bool) (*Client, error) {
	client := &Client{
		ID:       entities.NewID(),
		Code:     code,
		Name:     name,
		IsActive: isActive,
	}

	return client, nil
}
