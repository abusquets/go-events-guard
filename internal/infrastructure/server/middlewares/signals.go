package middlewares

import (
	"eventsguard/internal/infrastructure/signals"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

func NewSignalsMiddleware(
	api huma.API,
	asignalsBus signals.SignalsBus,
) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		defer func() {
			fmt.Println("SignalsMiddleware", ctx.URL().Path)
			asignalsBus.ProcessQueue()
		}()

		next(ctx)
	}
}
