package usecases

import (
	"context"

	"eventsguard/internal/auth/dtos"

	"eventsguard/internal/auth/domain/entities"

	"eventsguard/internal/app/errors"
)

type RefreshTokenUseCase interface {
	Execute(ctx context.Context, data dtos.RefreshTokenInputDTO) (*entities.Token, *errors.AppError)
}
