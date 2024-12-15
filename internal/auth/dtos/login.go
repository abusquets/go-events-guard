package dtos

import "eventsguard/internal/auth/domain/entities"

type LoginInputDTO struct {
	Username string `json:"username" minLength:"1"`
	Password string `json:"password" minLength:"1"`
}

type LoginRequest struct {
	Body LoginInputDTO `json:"body" validate:"required"`
}

type LoginResponse struct {
	Body *entities.Token
}
