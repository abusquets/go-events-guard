package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/constants"
	"eventsguard/internal/auth/domain/entities"
	auth_usecases_ports "eventsguard/internal/auth/domain/ports/usecases"
	"eventsguard/internal/auth/dtos"
	context_keys "eventsguard/internal/infrastructure/config/context"
	"eventsguard/internal/infrastructure/server/server_errors"
)

type LoginRouterHandler struct{}

func NewLoginRouterHandler(
	api huma.API,
	loginUseCase auth_usecases_ports.LoginUseCase,
	logoutUseCase auth_usecases_ports.LogoutUseCase,
) *LoginRouterHandler {
	ret := LoginRouterHandler{}

	huma.Register(api, huma.Operation{
		OperationID:   "login",
		Method:        http.MethodPost,
		Path:          "/auth/login",
		Description:   "Login",
		Tags:          []string{"Auth", "Login"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *dtos.LoginRequest) (*dtos.LoginResponse, error) {
		token, err := loginUseCase.Execute(
			ctx, dtos.LoginInputDTO{
				Username: input.Body.Username,
				Password: input.Body.Password,
			}, constants.DeviceApi,
		)

		if err != nil {
			return nil, server_errors.AppErrorToHumaError(*err)
		}

		resp := &dtos.LoginResponse{
			Body: token,
		}
		return resp, nil

	})

	huma.Register(api, huma.Operation{
		OperationID:   "logout",
		Method:        http.MethodPost,
		Path:          "/auth/logout",
		Description:   "Logout",
		Tags:          []string{"Auth", "Login"},
		DefaultStatus: http.StatusNoContent,
		Security: []map[string][]string{
			{"TokenAuth": {"api"}},
		},
	}, func(ctx context.Context, input *struct{}) (*struct{}, error) {

		token := ctx.Value(context_keys.TokenContextKey).(*entities.Token)
		done, err := logoutUseCase.Execute(ctx, token)

		if err != nil {
			return nil, server_errors.AppErrorToHumaError(*err)
		}

		if !done {
			return nil, server_errors.AppErrorToHumaError(
				*errors.NewNotFoundError("Error deleting token"),
			)

		}

		return nil, nil

	})

	return &ret
}
