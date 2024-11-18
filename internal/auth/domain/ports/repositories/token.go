package repositories

import (
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/constants"
	"eventsguard/internal/auth/domain/entities"
	"eventsguard/internal/auth/dtos"
)

type TokenRepository interface {
	Save(token dtos.CreateRawTokenDTO) *errors.AppError
	FindByToken(tokenKey string) (*entities.RawToken, *errors.AppError)
	FindByUserID(userID string) (map[constants.TokenDevice]*entities.RawToken, *errors.AppError)
	FindByUserIDAndDevice(userID string, device constants.TokenDevice) (*entities.RawToken, *errors.AppError)
	DeleteByToken(tokenKey string) (bool, *errors.AppError)
	DeleteByUserID(userID string) *errors.AppError
}
