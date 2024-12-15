package services

import (
	auth_service_ports "eventsguard/internal/auth/domain/ports/services"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewTokenService, fx.As(new(auth_service_ports.TokenService))),
	),
)
