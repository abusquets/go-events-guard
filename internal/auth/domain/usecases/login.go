package usecases

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/constants"
	"eventsguard/internal/auth/domain/entities"
	auth_services_ports "eventsguard/internal/auth/domain/ports/services"
	auth_usecases_ports "eventsguard/internal/auth/domain/ports/usecases"
	"eventsguard/internal/auth/dtos"
	core_services_ports "eventsguard/internal/core/domain/ports/services"
	"eventsguard/internal/infrastructure/mylog"
	"net/http"
)

const invalidCredentialsError = "Invalid username/password or user doesn't exist"

type loginUseCase struct {
	userService  core_services_ports.UserService
	tokenService auth_services_ports.TokenService
	logger       mylog.Logger
}

func NewLoginUseCase(
	userService core_services_ports.UserService,
	tokenService auth_services_ports.TokenService,
) auth_usecases_ports.LoginUseCase {
	return loginUseCase{
		userService:  userService,
		tokenService: tokenService,
		logger:       mylog.GetLogger(),
	}
}

func (u loginUseCase) Execute(
	ctx context.Context,
	data dtos.LoginInputDTO,
	device constants.TokenDevice,
) (token *entities.Token, error *errors.AppError) {
	user, error := u.userService.GetUserByEmail(ctx, data.Username)

	if error != nil || user == nil {
		if error.Code == http.StatusNotFound {
			return nil, errors.NewAuthError(
				invalidCredentialsError,
			)
		}
		return nil, error
	}

	if !user.VerifyPassword(data.Password) {
		u.logger.Debug("Invalid password")
		return nil, errors.NewAuthError(
			invalidCredentialsError,
		)
	}

	if !user.IsActive {
		u.logger.Debug("User is not active")
		return nil, errors.NewAuthError(
			invalidCredentialsError,
		)
	}

	fakeUser := entities.FakeUser{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Email,
		IsAdmin:   user.IsAdmin,
	}
	token = new(entities.Token)
	token, err := u.tokenService.CreateByUser(fakeUser, device, true, nil)
	if err != nil {
		return nil, errors.NewUnexpectedError(err.Error())
	}
	if token == nil {
		return nil, errors.NewUnexpectedError("Error creating token")
	}

	return token, nil
}
