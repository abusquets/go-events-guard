package inbound_events

import (
	"eventsguard/internal/inbound_events/adapters/api"
	// "eventsguard/internal/auth/adapters/spi/repositories"
	// "eventsguard/internal/auth/domain/services"
	"eventsguard/internal/inbound_events/domain/usecases"

	"go.uber.org/fx"
)

var Module = fx.Module("inbound_events",
	// repositories.Module,
	// services.Module,
	usecases.Module,
	api.Module,
)
