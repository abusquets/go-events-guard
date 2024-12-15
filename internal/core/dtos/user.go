package dtos

import (
	"eventsguard/internal/core/domain/entities"
	utils_entities "eventsguard/internal/utils/entities"

	"eventsguard/internal/utils/dtos/pagination"

	"github.com/guregu/null"
)

type CreateUserInput struct {
	Email     string  `json:"email" validate:"required,email"`
	FirstName string  `json:"first_name" validate:"required"`
	LastName  *string `json:"last_name,omitempty"`
	Password  string  `json:"password" validate:"required"`
	IsActive  bool    `json:"is_active,omitempty"`
}

type CreateUserRequest struct {
	Body CreateUserInput `json:"body" validate:"required"`
}

type CreateUserResponse struct {
	Body entities.User
}

type ListUserResponse struct {
	Body *pagination.PaginatedResult[entities.User]
}

type UserDetailResponse struct {
	Body entities.User
}

type UpdatePartialUserInput struct {
	FirstName null.String         `json:"first_name,omitempty"`
	LastName  null.String         `json:"last_name,omitempty"`
	IsActive  *bool               `json:"is_active,omitempty"`
	Clients   []utils_entities.ID `json:"clients,omitempty"`
}

type UpdatePartialAdminUserInput struct {
	FirstName null.String         `json:"first_name,omitempty"`
	LastName  null.String         `json:"last_name,omitempty"`
	IsActive  *bool               `json:"is_active,omitempty"`
	IsAdmin   *bool               `json:"is_admin,omitempty"`
	Clients   []utils_entities.ID `json:"clients,omitempty"`
}

type UpdatePartialUserRequest struct {
	Id   string                 `path:"id" maxLength:"24" example:"671495c80d8e9fd6b3256340" doc:"User ID"`
	Body UpdatePartialUserInput `json:"body" validate:"required"`
}
