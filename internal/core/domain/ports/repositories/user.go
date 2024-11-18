package repositories

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/core/domain/entities"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/utils/dtos/pagination"
)

type UserQuery struct {
	Page     *int
	PageSize *int
	Search   *string
}

type UserRepository interface {
	GetByID(ctx context.Context, ID string) (*entities.User, *errors.AppError)
	GetByEmail(ctx context.Context, Email string) (*entities.User, *errors.AppError)
	Create(ctx context.Context, userData dtos.CreateUserInput) (user *entities.User, error *errors.AppError)
	List(ctx context.Context, query UserQuery) (*pagination.PaginatedResult[entities.User], *errors.AppError)
	UpdatePartialUser(
		ctx context.Context,
		ID string,
		userData dtos.UpdatePartialAdminUserInput,
	) (*entities.User, *errors.AppError)
}
