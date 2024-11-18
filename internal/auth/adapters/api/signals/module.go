package signals

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewAuthUserDelayedSignalsHandler,
	),
	fx.Invoke(
		func(
			authUserDelayedSignalsHandler *AuthUserDelayedSignalsHandler,
		) {
		},
	),
)
