package server

import (
	"context"
	"net/http"

	"go.uber.org/fx"
)

type WebServer interface {
	Start()
	Shutdown(ctx context.Context) error
	GetMux() *http.ServeMux
	RegisterHooks(lifecycle fx.Lifecycle)
}

// type ApiServer interface {
// 	Start()
// 	// RegisterHooks(lifecycle fx.Lifecycle)
// }

type Worker interface {
	Start(ctx context.Context)
	RegisterHooks(lifecycle fx.Lifecycle)
}
