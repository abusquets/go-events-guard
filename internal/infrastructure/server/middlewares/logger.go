package middlewares

import (
	"eventsguard/internal/infrastructure/mylog"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

func NewLoggerMiddleware(
	api huma.API,
) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		logger := mylog.GetLogger()

		start := time.Now()
		next(ctx)
		elapsed := time.Since(start)
		logger.With("uri", ctx.Operation().Path).
			With("query", ctx.Operation().Parameters).
			With("method", ctx.Operation().Method).
			With("status", ctx.Status()).
			With("elapsed", elapsed.Round(time.Millisecond).String()).
			Info("RequestLog")
	}
}
