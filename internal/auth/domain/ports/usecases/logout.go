package usecases

import (
	"context"

	"eventsguard/internal/auth/domain/entities"

	"eventsguard/internal/app/errors"
)

//go:generate mockgen -source=logout.go -destination=mock/logout.go

type LogoutUseCase interface {
	Execute(ctx context.Context, token *entities.Token) (bool, *errors.AppError)
}
