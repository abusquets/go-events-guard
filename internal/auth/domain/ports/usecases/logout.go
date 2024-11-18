package usecases

import (
	"context"

	"eventsguard/internal/auth/domain/entities"

	"eventsguard/internal/app/errors"
)

type LogoutUseCase interface {
	Execute(ctx context.Context, token *entities.Token) (bool, *errors.AppError)
}
