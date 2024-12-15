// internal/app/module.go.
package app

import (
	"eventsguard/internal/auth"
	"eventsguard/internal/core"
	"eventsguard/internal/inbound_events"
	"eventsguard/internal/infrastructure"
	"eventsguard/internal/infrastructure/config"
	"eventsguard/internal/infrastructure/migrations"

	"go.uber.org/fx"
)

var Module = fx.Module("app",
	config.Module,
	infrastructure.Module,
	core.Module,
	auth.Module,
	inbound_events.Module,
	migrations.Module,
)
