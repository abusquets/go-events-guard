package dtos

import (
	"eventsguard/internal/core/domain/entities"
)

type CreatePersonInput struct {
	FirstName string  `json:"first_name" validate:"required"`
	LastName  *string `json:"last_name,omitempty"` // Optional field`
}

type CreatePersonRequest struct {
	Body CreatePersonInput `json:"body" validate:"required"`
}

type CreatePersonResponse struct {
	Body entities.Person
}

type ListPersonResponse struct {
	Body []entities.Person
}

type PersonDetailResponse struct {
	Body entities.Person
}
