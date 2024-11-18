// migrations/module.go
package migrations

import (
	"go.uber.org/fx"
)

// MigrationsModule exports the migration commands and dependencies
var Module = fx.Module(
	"migrations",
	fx.Provide(
		NewMigrateCommands,
	),
)
