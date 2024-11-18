package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"eventsguard/internal/auth/domain/entities"
	core_repositories_ports "eventsguard/internal/core/domain/ports/repositories"
	core_service_ports "eventsguard/internal/core/domain/ports/services"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/infrastructure/config"
	ctx_config "eventsguard/internal/infrastructure/config/context"
	"eventsguard/internal/infrastructure/server/server_errors"
)

type UserRouterHandler struct {
}

func NewUserRouterHandler(
	api huma.API,
	userService core_service_ports.UserService,
	cfg *config.AppConfig,
) *UserRouterHandler {

	ret := UserRouterHandler{}

	registerAddUser(api, userService)
	registerListUser(api, userService)
	registerUserDetail(api, userService)
	registerUserPartialUpdate(api, userService)

	return &ret
}

func registerAddUser(api huma.API, userService core_service_ports.UserService) {
	huma.Register(api, huma.Operation{
		OperationID:   "add-user",
		Method:        http.MethodPost,
		Path:          "/user",
		Description:   "Add a User",
		Tags:          []string{"Core", "User"},
		DefaultStatus: http.StatusCreated,
		Security: []map[string][]string{
			{"TokenAuth": {"api"}},
		},
	}, func(ctx context.Context, input *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
		newUser, err := userService.CreateUser(ctx, input.Body)

		if err != nil {
			return nil, server_errors.AppErrorToHumaError(*err)
		}

		resp := &dtos.CreateUserResponse{
			Body: *newUser,
		}
		return resp, nil
	})
}

func registerListUser(api huma.API, userService core_service_ports.UserService) {
	huma.Register(api, huma.Operation{
		OperationID:   "list-user",
		Method:        http.MethodGet,
		Path:          "/user",
		Description:   "Get the list of users",
		Tags:          []string{"Core", "User"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"TokenAuth": {"api"}},
		},
	}, func(ctx context.Context, input *struct {
		Page     int    `query:"page" doc:"Page number"`
		PageSize int    `query:"pageSize" doc:"Number of items per page"`
		Search   string `query:"search" doc:"Search query"`
	}) (*dtos.ListUserResponse, error) {
		query := core_repositories_ports.UserQuery{
			Page:     &input.Page,
			PageSize: &input.PageSize,
			Search:   &input.Search,
		}
		users, err := userService.ListUsers(ctx, query)

		if err != nil {
			return nil, server_errors.AppErrorToHumaError(*err)
		}

		resp := &dtos.ListUserResponse{
			Body: users,
		}
		return resp, nil
	})
}

func registerUserDetail(api huma.API, userService core_service_ports.UserService) {
	huma.Register(api, huma.Operation{
		OperationID:   "user-detail",
		Method:        http.MethodGet,
		Path:          "/user/{id}",
		Description:   "Get a user by ID",
		Tags:          []string{"Core", "User"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"TokenAuth": {"api"}},
		},
	}, func(ctx context.Context, input *struct {
		Id string `path:"id" maxLength:"24" example:"671495c80d8e9fd6b3256340" doc:"User ID"`
	}) (*dtos.UserDetailResponse, error) {
		user, err := userService.GetUserByID(ctx, input.Id)

		if err != nil {
			return nil, server_errors.AppErrorToHumaError(*err)
		}

		resp := &dtos.UserDetailResponse{
			Body: *user,
		}
		return resp, nil
	})
}

func registerUserPartialUpdate(api huma.API, userService core_service_ports.UserService) {
	huma.Register(api, huma.Operation{
		OperationID:   "user-partial-update",
		Method:        http.MethodPatch,
		Path:          "/user/{id}",
		Description:   "Update partial user by ID",
		Tags:          []string{"Core", "User"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"TokenAuth": {"api"}},
		},
	}, func(ctx context.Context, input *dtos.UpdatePartialUserRequest) (*dtos.UserDetailResponse, error) {

		appContext := context.Background()
		tokenValue := ctx.Value(ctx_config.TokenContextKey)
		token, _ := tokenValue.(*entities.Token)
		appContext = context.WithValue(appContext, ctx_config.UserContextKey, token.User)

		inputData := dtos.UpdatePartialAdminUserInput{}
		if input.Body.FirstName.Valid {
			inputData.FirstName = input.Body.FirstName
		}
		if input.Body.LastName.Valid {
			inputData.LastName = input.Body.LastName
		}
		if input.Body.IsActive != nil {
			inputData.IsActive = input.Body.IsActive
		}
		if input.Body.Clients != nil {
			inputData.Clients = input.Body.Clients
		}

		currentUser, err := userService.UpdatePartialUser(appContext, input.Id, inputData)

		if err != nil {
			return nil, server_errors.AppErrorToHumaError(*err)
		}

		resp := &dtos.UserDetailResponse{
			Body: *currentUser,
		}
		return resp, nil
	})
}
