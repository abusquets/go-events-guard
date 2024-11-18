package services

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/core/domain/entities"
	"eventsguard/internal/core/domain/ports/repositories"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/utils/dtos/pagination"
)

type UserService interface {
	GetUserByID(ctx context.Context, ID string) (*entities.User, *errors.AppError)
	GetUserByEmail(ctx context.Context, Email string) (*entities.User, *errors.AppError)
	CreateUser(ctx context.Context, userData dtos.CreateUserInput) (user *entities.User, error *errors.AppError)
	ListUsers(ctx context.Context, query repositories.UserQuery) (*pagination.PaginatedResult[entities.User], *errors.AppError)
	UpdatePartialUser(
		ctx context.Context,
		ID string,
		userData dtos.UpdatePartialAdminUserInput,
	) (*entities.User, *errors.AppError)
}
