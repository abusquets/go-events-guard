package core

import (
	"eventsguard/internal/core/adapters/api"
	"eventsguard/internal/core/adapters/spi/repositories"
	"eventsguard/internal/core/domain/services"

	"go.uber.org/fx"
)

var Module = fx.Module("core",
	repositories.Module,
	services.Module,
	api.Module,
)
