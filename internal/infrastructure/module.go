// internal/infrastructure/module.go
package infrastructure

import (
	"eventsguard/internal/infrastructure/database"
	"eventsguard/internal/infrastructure/redis"

	"eventsguard/internal/infrastructure/server"
	"eventsguard/internal/infrastructure/signals"

	"go.uber.org/fx"
)

var Module = fx.Options(
	database.Module,
	server.Module,
	redis.Module,
	signals.Module,
)
