package usecases

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/domain/entities"
	auth_services_ports "eventsguard/internal/auth/domain/ports/services"
	auth_usecases_ports "eventsguard/internal/auth/domain/ports/usecases"
)

type logoutUseCase struct {
	tokenService auth_services_ports.TokenService
}

func NewLogoutUseCase(
	tokenService auth_services_ports.TokenService,
) auth_usecases_ports.LogoutUseCase {
	return logoutUseCase{
		tokenService: tokenService,
	}
}

func (u logoutUseCase) Execute(
	ctx context.Context,
	token *entities.Token,
) (bool, *errors.AppError) {
	return u.tokenService.DeleteByToken(token.Token)
}
