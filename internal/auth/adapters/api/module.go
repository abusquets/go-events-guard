package api

import (
	"eventsguard/internal/auth/adapters/api/http"
	"eventsguard/internal/auth/adapters/api/signals"

	"go.uber.org/fx"
)

var Module = fx.Options(
	http.Module,
	signals.Module,
)
