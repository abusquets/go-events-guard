package http

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewLoginRouterHandler,
		NewTokenRouterHandler,
	),
	fx.Invoke(
		func(
			loginRouterHandler *LoginRouterHandler,
			tokenRouterHandler *TokenRouterHandler,
		) {
			// Do nothing
		},
	),
)
