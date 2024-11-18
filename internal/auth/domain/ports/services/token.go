package services

import (
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/constants"
	"eventsguard/internal/auth/domain/entities"
)

type TokenService interface {
	CreateByUser(user entities.FakeUser, device constants.TokenDevice, expirable bool, renew *bool) (*entities.Token, *errors.AppError)
	FindByTokenKey(tokenKey string) (*entities.Token, *errors.AppError)
	DeleteByToken(tokenKey string) (bool, *errors.AppError)
	DeleteByUserID(userID string, devices []constants.TokenDevice) *errors.AppError
	RenewToken(token entities.Token) *errors.AppError
	IsTokenExpired(token *entities.Token) bool
	FindByUserID(userID string) (map[constants.TokenDevice]*entities.Token, *errors.AppError)
}
