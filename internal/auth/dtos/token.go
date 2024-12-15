package dtos

import (
	"eventsguard/internal/auth/constants"
	"eventsguard/internal/auth/domain/entities"
	"time"
)

type CreateRawTokenDTO struct {
	Device    constants.TokenDevice `json:"device"`
	Token     string                `json:"token"`
	Payload   string                `json:"payload"`
	UserID    string                `json:"user_id"`
	ExpiresAt *time.Time            `json:"expires_at"`
	CreatedAt time.Time             `json:"created_at"`
	Expiracy  *int                  `json:"expiracy"`
}

type RefreshTokenInputDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	Body RefreshTokenInputDTO `json:"body" validate:"required"`
}

type RefreshTokenResponse struct {
	Body *entities.Token
}
