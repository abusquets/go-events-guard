package usecases

import (
	"context"

	"eventsguard/internal/auth/constants"

	"eventsguard/internal/auth/dtos"

	"eventsguard/internal/auth/domain/entities"

	"eventsguard/internal/app/errors"
)

type LoginUseCase interface {
	Execute(ctx context.Context, data dtos.LoginInputDTO, device constants.TokenDevice) (*entities.Token, *errors.AppError)
}
