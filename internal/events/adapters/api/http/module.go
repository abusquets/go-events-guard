package http

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewEventRouterHandler,
	),
	fx.Invoke(
		func(
			eventRouterHandler *EventRouterHandler,
		) {
		},
	),
)
