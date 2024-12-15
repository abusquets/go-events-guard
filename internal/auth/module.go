package auth

import (
	"eventsguard/internal/auth/adapters/api"
	"eventsguard/internal/auth/adapters/spi/repositories"
	"eventsguard/internal/auth/domain/services"
	"eventsguard/internal/auth/domain/usecases"

	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	repositories.Module,
	services.Module,
	usecases.Module,
	api.Module,
)
