package services

import (
	core_service_ports "eventsguard/internal/core/domain/ports/services"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewUserService, fx.As(new(core_service_ports.UserService))),
		fx.Annotate(NewClientService, fx.As(new(core_service_ports.ClientService))),
	),
)
