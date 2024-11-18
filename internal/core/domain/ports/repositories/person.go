package repositories

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/core/domain/entities"
	"eventsguard/internal/core/dtos"
)

type PersonRepository interface {
	GetByID(ctx context.Context, ID string) (*entities.Person, *errors.AppError)
	Create(ctx context.Context, userData dtos.CreatePersonInput) (user *entities.Person, error *errors.AppError)
	List(ctx context.Context) (*[]entities.Person, *errors.AppError)
}
