package services

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/core/domain/entities"
	core_repositories_ports "eventsguard/internal/core/domain/ports/repositories"
	core_service_ports "eventsguard/internal/core/domain/ports/services"
	"eventsguard/internal/core/dtos"
)

type personService struct {
	personRepository core_repositories_ports.PersonRepository
}

func NewPersonService(personRepository core_repositories_ports.PersonRepository) core_service_ports.PersonService {
	return &personService{
		personRepository: personRepository,
	}
}

// CreatePerson now accepts context.Context and passes it to the repository
func (u *personService) CreatePerson(ctx context.Context, personData dtos.CreatePersonInput) (*entities.Person, *errors.AppError) {
	return u.personRepository.Create(ctx, personData)
}

// GetPersonByID now accepts context.Context and passes it to the repository
func (u *personService) GetPersonByID(ctx context.Context, ID string) (*entities.Person, *errors.AppError) {
	return u.personRepository.GetByID(ctx, ID)
}

// ListPersons now accepts context.Context and passes it to the repository
func (u *personService) ListPersons(ctx context.Context) (*[]entities.Person, *errors.AppError) {
	return u.personRepository.List(ctx)
}
