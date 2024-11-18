package api

import (
	"eventsguard/internal/core/adapters/api/cli"
	"eventsguard/internal/core/adapters/api/http"

	"go.uber.org/fx"
)

var Module = fx.Options(
	http.Module,
	cli.Module,
)
