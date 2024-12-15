package redis

import (
	auth_repository_ports "eventsguard/internal/auth/domain/ports/repositories"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewRedisTokenRepository, fx.As(new(auth_repository_ports.TokenRepository))),
	),
)
