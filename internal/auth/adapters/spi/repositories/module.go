package repositories

import (
	"eventsguard/internal/auth/adapters/spi/repositories/redis"

	"go.uber.org/fx"
)

var Module = fx.Options(
	redis.Module,
)
