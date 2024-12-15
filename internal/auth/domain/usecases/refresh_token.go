package usecases

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/domain/entities"
	auth_services_ports "eventsguard/internal/auth/domain/ports/services"
	auth_usecases_ports "eventsguard/internal/auth/domain/ports/usecases"
	"eventsguard/internal/auth/dtos"
	core_services_ports "eventsguard/internal/core/domain/ports/services"
	"eventsguard/internal/infrastructure/mylog"
)

type refreshTokenUseCase struct {
	userService  core_services_ports.UserService
	tokenService auth_services_ports.TokenService
	logger       mylog.Logger
}

func NewRefreshTokenUseCase(
	userService core_services_ports.UserService,
	tokenService auth_services_ports.TokenService,
) auth_usecases_ports.RefreshTokenUseCase {
	return refreshTokenUseCase{
		userService:  userService,
		tokenService: tokenService,
		logger:       mylog.GetLogger(),
	}
}

func (u refreshTokenUseCase) Execute(
	ctx context.Context,
	data dtos.RefreshTokenInputDTO,
) (token *entities.Token, error *errors.AppError) {

	// token, error := u.tokenService.RefreshToken(data.AccessToken, data.RefreshToken)

	// if error != nil {
	// 	if error.Code == http.StatusNotFound {
	// 		u.logger.Debug("User not found")
	// 		return nil, errors.NewValidationError(
	// 			"Invalid username/password or user doesn't exist",
	// 		)
	// 	}
	// 	return nil, error
	// }

	// if user == nil {
	// 	u.logger.Debug("User not found")
	// 	return nil, errors.NewValidationError(
	// 		"Invalid username/password or user doesn't exist",
	// 	)
	// }

	// if !user.VerifyPassword(data.Password) {
	// 	u.logger.Debug("Invalid password")
	// 	return nil, errors.NewValidationError(
	// 		"Invalid username/password or user doesn't exist",
	// 	)
	// }

	// if !user.IsActive {
	// 	u.logger.Debug("User is not active")
	// 	return nil, errors.NewValidationError(
	// 		"Invalid username/password or user doesn't exist",
	// 	)
	// }

	// fakeUser := entities.FakeUser{
	// 	ID:        user.ID.String(),
	// 	FirstName: user.FirstName,
	// 	LastName:  user.LastName,
	// 	Username:  user.Email,
	// 	IsAdmin:   user.IsAdmin,
	// }
	// token = new(entities.Token)
	// token, err := u.tokenService.CreateByUser(fakeUser, device, true, nil)
	// if err != nil {
	// 	return nil, errors.NewUnexpectedError(err.Error())
	// }
	// if token == nil {
	// 	return nil, errors.NewUnexpectedError("Error creating token")
	// }

	return token, nil
}
