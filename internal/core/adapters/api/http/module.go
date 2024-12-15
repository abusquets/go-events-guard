package http

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewUserRouterHandler,
		NewClientRouterHandler,
	),
	fx.Invoke(
		func(
			userRouterHandler *UserRouterHandler,
			clientRouterHandler *ClientRouterHandler,
		) {
		},
	),
)
