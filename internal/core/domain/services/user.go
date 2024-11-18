package services

import (
	"context"
	"eventsguard/internal/app/errors"
	auth_entities "eventsguard/internal/auth/domain/entities"
	"eventsguard/internal/core/domain/entities"
	"eventsguard/internal/core/domain/ports/repositories"
	core_repository_ports "eventsguard/internal/core/domain/ports/repositories"
	core_service_ports "eventsguard/internal/core/domain/ports/services"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/infrastructure/config"
	ctx_config "eventsguard/internal/infrastructure/config/context"
	"eventsguard/internal/infrastructure/mylog"
	"eventsguard/internal/infrastructure/signals"
	"eventsguard/internal/utils/dtos/pagination"
)

type userService struct {
	userRepository core_repository_ports.UserRepository
	logger         mylog.Logger
	asignalsBus    signals.SignalsBus
}

func NewUserService(
	cfg *config.AppConfig,
	userRepository core_repository_ports.UserRepository,
	asignalsBus signals.SignalsBus,
) core_service_ports.UserService {
	return userService{
		userRepository: userRepository,
		logger:         mylog.GetLogger(),
		asignalsBus:    asignalsBus,
	}
}

func (us userService) CreateUser(ctx context.Context, userData dtos.CreateUserInput) (*entities.User, *errors.AppError) {
	return us.userRepository.Create(ctx, userData)
}

func (us userService) GetUserByEmail(ctx context.Context, Email string) (*entities.User, *errors.AppError) {
	return us.userRepository.GetByEmail(ctx, Email)
}

func (us userService) GetUserByID(ctx context.Context, ID string) (*entities.User, *errors.AppError) {
	return us.userRepository.GetByID(ctx, ID)
}

func (us userService) ListUsers(ctx context.Context, query repositories.UserQuery) (*pagination.PaginatedResult[entities.User], *errors.AppError) {
	return us.userRepository.List(ctx, query)
}

func (us userService) UpdatePartialUser(
	ctx context.Context,
	ID string,
	userData dtos.UpdatePartialAdminUserInput,
) (*entities.User, *errors.AppError) {
	responsible, ok := ctx.Value(ctx_config.UserContextKey).(auth_entities.FakeUser)
	if ok {
		if !responsible.IsAdmin && responsible.ID != ID {
			return nil, errors.NewPermissionDeniedError(
				"Forbidden: user cannot update other users",
			)
		}
	}
	user, error := us.userRepository.UpdatePartialUser(ctx, ID, userData)
	if error == nil {
		us.asignalsBus.AfterTransaction("user:updated", user.ID)
	}
	return user, error
}
