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
