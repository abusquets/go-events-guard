package repositories

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/core/domain/entities"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/utils/dtos/pagination"
)

type ClientQuery struct {
	Page     *int
	PageSize *int
	Search   *string
}

type ClientRepository interface {
	GetByID(ctx context.Context, ID string) (*entities.Client, *errors.AppError)
	GetByEmail(ctx context.Context, Email string) (*entities.Client, *errors.AppError)
	Create(ctx context.Context, clientData dtos.CreateClientInput) (client *entities.Client, error *errors.AppError)
	List(ctx context.Context, query ClientQuery) (*pagination.PaginatedResult[entities.Client], *errors.AppError)
	UpdatePartialClient(
		ctx context.Context,
		ID string,
		clientData dtos.UpdatePartialClientInput,
	) (*entities.Client, *errors.AppError)
}
