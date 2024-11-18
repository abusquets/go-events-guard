package services

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/core/domain/entities"
	"eventsguard/internal/core/dtos"
)

type PersonService interface {
	GetPersonByID(ctx context.Context, ID string) (*entities.Person, *errors.AppError)
	CreatePerson(ctx context.Context, userData dtos.CreatePersonInput) (user *entities.Person, error *errors.AppError)
	ListPersons(ctx context.Context) (*[]entities.Person, *errors.AppError)
}
