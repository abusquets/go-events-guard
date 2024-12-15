package entities

import (
	"eventsguard/internal/auth/constants"
	"time"
)

type FakeUser struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Username  string  `json:"username"`
	IsAdmin   bool    `json:"is_admin"`
}

type Token struct {
	Device       constants.TokenDevice `json:"device"`
	Token        string                `json:"token"`
	RefreshToken string                `json:"refresh_token"`
	User         FakeUser              `json:"user"`
	UserID       string                `json:"user_id"`
	ExpiresAt    *time.Time            `json:"expires_at,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
	Expiracy     *int                  `json:"expiracy,omitempty"`
}

type RawToken struct {
	Device       constants.TokenDevice `json:"device"`
	Token        string                `json:"token"`
	RefreshToken string                `json:"refresh_token"`
	Payload      string                `json:"payload"`
	UserID       string                `json:"user_id"`
	ExpiresAt    *time.Time            `json:"expires_at,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
	Expiracy     *int                  `json:"expiracy,omitempty"`
}
