package mongodb

import (
	core_repository_ports "eventsguard/internal/core/domain/ports/repositories"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewUserRepository, fx.As(new(core_repository_ports.UserRepository))),
		fx.Annotate(NewClientRepository, fx.As(new(core_repository_ports.ClientRepository))),
		fx.Annotate(NewPersonRepository, fx.As(new(core_repository_ports.PersonRepository))),
	),
)
