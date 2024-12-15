package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	auth_usecases_ports "eventsguard/internal/auth/domain/ports/usecases"
	"eventsguard/internal/auth/dtos"
	"eventsguard/internal/infrastructure/server/server_errors"
)

type TokenRouterHandler struct{}

func NewTokenRouterHandler(
	api huma.API,
	refreshTokenUseCase auth_usecases_ports.RefreshTokenUseCase,
) *TokenRouterHandler {
	ret := TokenRouterHandler{}

	huma.Register(api, huma.Operation{
		OperationID:   "refresh-token",
		Method:        http.MethodPost,
		Path:          "/auth/refresh-token",
		Description:   "Refresh token",
		Tags:          []string{"Auth", "Token"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *dtos.RefreshTokenRequest) (*dtos.RefreshTokenResponse, error) {
		token, err := refreshTokenUseCase.Execute(
			ctx, input.Body,
		)

		if err != nil {
			return nil, server_errors.AppErrorToHumaError(*err)
		}

		resp := &dtos.RefreshTokenResponse{
			Body: token,
		}
		return resp, nil

	})

	return &ret
}
