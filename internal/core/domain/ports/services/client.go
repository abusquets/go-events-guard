package services

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/core/domain/entities"
	"eventsguard/internal/core/domain/ports/repositories"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/utils/dtos/pagination"
)

type ClientService interface {
	GetClientByID(ctx context.Context, ID string) (*entities.Client, *errors.AppError)
	GetClientByEmail(ctx context.Context, Email string) (*entities.Client, *errors.AppError)
	CreateClient(ctx context.Context, clientData dtos.CreateClientInput) (client *entities.Client, error *errors.AppError)
	ListClients(ctx context.Context, query repositories.ClientQuery) (*pagination.PaginatedResult[entities.Client], *errors.AppError)
	UpdatePartialClient(
		ctx context.Context,
		ID string,
		clientData dtos.UpdatePartialClientInput,
	) (*entities.Client, *errors.AppError)
}
