package repositories

import (
	"eventsguard/internal/core/adapters/spi/repositories/mongodb"

	"go.uber.org/fx"
)

var Module = fx.Options(
	mongodb.Module,
)
