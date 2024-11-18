package dtos

import "eventsguard/internal/auth/domain/entities"

type LoginInputDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Body LoginInputDTO `json:"body" validate:"required"`
}

type LoginResponse struct {
	Body *entities.Token
}
