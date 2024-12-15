package http

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewInboundEventRouterHandler,
	),
	fx.Invoke(
		func(
			eventRouterHandler *InboundEventRouterHandler,
		) {
		},
	),
)
