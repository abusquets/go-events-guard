package usecases

import (
	"go.uber.org/fx"
)

// var Module = fx.Options(
// 	fx.Provide(
// 		fx.Annotate(NewLoginUseCase, fx.As(new(auth_usecases_ports.LoginUseCase))),
// 		fx.Annotate(NewLogoutUseCase, fx.As(new(auth_usecases_ports.LogoutUseCase))),
// 		fx.Annotate(NewRefreshTokenUseCase, fx.As(new(auth_usecases_ports.RefreshTokenUseCase))),
// 	),
// )

var Module = fx.Options(
	fx.Provide(
		NewLoginUseCase,
		NewLogoutUseCase,
		NewRefreshTokenUseCase,
	),
)
