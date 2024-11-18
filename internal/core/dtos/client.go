package dtos

import (
	"eventsguard/internal/core/domain/entities"

	"eventsguard/internal/utils/dtos/pagination"

	"github.com/guregu/null"
)

type CreateClientInput struct {
	Code     string `json:"code" validate:"required"`
	Name     string `json:"name" validate:"required"`
	IsActive bool   `json:"is_active,omitempty"`
}

type CreateClientRequest struct {
	Body CreateClientInput `json:"body" validate:"required"`
}

type CreateClientResponse struct {
	Body entities.Client
}

type ListClientResponse struct {
	Body *pagination.PaginatedResult[entities.Client]
}

type ClientDetailResponse struct {
	Body entities.Client
}

type UpdatePartialClientInput struct {
	Name     null.String `json:"name,omitempty"`
	IsActive *bool       `json:"is_active,omitempty"`
}

type UpdatePartialClientRequest struct {
	Id   string                   `path:"id" maxLength:"24" example:"671495c80d8e9fd6b3256340" doc:"Client ID"`
	Body UpdatePartialClientInput `json:"body" validate:"required"`
}
