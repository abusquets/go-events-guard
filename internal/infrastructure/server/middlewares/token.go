package middlewares

import (
	context_keys "eventsguard/internal/infrastructure/config/context"
	"net/http"
	"slices"
	"strings"

	auth_services "eventsguard/internal/auth/domain/ports/services"

	"github.com/danielgtaylor/huma/v2"
)

func NewAuthTokenMiddleware(
	api huma.API,
	tokenService auth_services.TokenService,
) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		var neededScopes []string
		isAuthorizationRequired := false
		for _, opScheme := range ctx.Operation().Security {
			var ok bool
			if neededScopes, ok = opScheme["TokenAuth"]; ok {
				isAuthorizationRequired = true
				break
			}
		}

		if !isAuthorizationRequired {
			next(ctx)
			return
		}

		tokenKey := strings.TrimPrefix(ctx.Header("Authorization"), "Bearer ")

		if len(tokenKey) == 0 {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}
		token, err := tokenService.FindByTokenKey(tokenKey)
		if err != nil || token == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if tokenService.IsTokenExpired(token) {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Expired token")
			return
		}

		if slices.Contains(neededScopes, string(token.Device)) {
			ctx.Context()
			newCtx := huma.WithValue(ctx, context_keys.TokenContextKey, token)
			next(newCtx)
			return
		}

		huma.WriteErr(api, ctx, http.StatusForbidden, "Forbidden")
	}
}
