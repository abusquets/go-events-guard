package api

import (
	"eventsguard/internal/inbound_events/adapters/api/http"
	// "eventsguard/internal/auth/adapters/api/signals"

	"go.uber.org/fx"
)

var Module = fx.Options(
	http.Module,
	// signals.Module,
)
