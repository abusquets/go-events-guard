package server

import (
	"fmt"

	auth_services "eventsguard/internal/auth/domain/ports/services"
	"eventsguard/internal/infrastructure/config"
	"eventsguard/internal/infrastructure/server/middlewares"
	"eventsguard/internal/infrastructure/signals"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type ApiWebServer struct {
	Api huma.API
}

func NewHumaServer(
	cfg *config.AppConfig,
	server WebServer,
	tokenService auth_services.TokenService,
	signalsBus signals.SignalsBus,
) huma.API {

	// Configuració de l'API Huma
	apiConfig := huma.DefaultConfig(cfg.ApiName, cfg.ApiVersion)
	apiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"TokenAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	url := fmt.Sprintf("%s/api", cfg.ServerUri)
	apiConfig.Servers = []*huma.Server{{URL: url}}

	// Creació de l'API Huma amb el prefix "/api"
	api := humago.NewWithPrefix(server.GetMux(), "/api", apiConfig)
	api.UseMiddleware(middlewares.NewLoggerMiddleware(api))
	api.UseMiddleware(middlewares.NewSignalsMiddleware(api, signalsBus))
	api.UseMiddleware(middlewares.NewAuthTokenMiddleware(api, tokenService))
	return api
}

// Start inicia el servidor HTTP.
// ApiWebServer modificat
// func (a *ApiWebServer) Start() {
// 	log.Println("Starting ApiWebServer...")
// 	a.server.Start() // Assegura't que el servidor web s'inicia en una goroutine separada
// }

// // RegisterHooks afegeix els hooks d'inici i aturada per ApiWebServer.
// func (a *ApiWebServer) RegisterHooks(lifecycle fx.Lifecycle) {
// 	lifecycle.Append(
// 		fx.Hook{
// 			OnStart: func(ctx context.Context) error {
// 				go a.Start()
// 				return nil
// 			},
// 			OnStop: func(ctx context.Context) error {
// 				log.Println("Shutting down ApiWebServer...")
// 				//return a.server.Shutdown(ctx)
// 				return nil
// 			},
// 		},
// 	)
// }
