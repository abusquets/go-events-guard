package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	core_service_ports "eventsguard/internal/core/domain/ports/services"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/infrastructure/config"
	"eventsguard/internal/infrastructure/server/server_errors"
)

type EventRouterHandler struct {
}

func NewEventRouterHandler(
	api huma.API,
	clientService core_service_ports.ClientService,
	cfg *config.AppConfig,
) *EventRouterHandler {
	ret := EventRouterHandler{}

	huma.Register(api, huma.Operation{
		OperationID:   "add-event",
		Method:        http.MethodPost,
		Path:          "/event",
		Description:   "Add an Event",
		Tags:          []string{"Event"},
		DefaultStatus: http.StatusCreated,
		Security: []map[string][]string{
			{"TokenAuth": {"api"}},
		},
	}, func(ctx context.Context, input *dtos.CreateClientRequest) (*dtos.CreateClientResponse, error) {
		newClient, err := clientService.CreateClient(ctx, input.Body)

		if err != nil {
			return nil, server_errors.AppErrorToHumaError(*err)
		}

		resp := &dtos.CreateClientResponse{
			Body: *newClient,
		}
		return resp, nil
	})

	// huma.Register(api, huma.Operation{
	// 	OperationID:   "list-client",
	// 	Method:        http.MethodGet,
	// 	Path:          "/client",
	// 	Description:   "Get the list of clients",
	// 	Tags:          []string{"Core", "Client"},
	// 	DefaultStatus: http.StatusOK,
	// 	Security: []map[string][]string{
	// 		{"TokenAuth": {"api"}},
	// 	},
	// }, func(ctx context.Context, input *struct {
	// 	Page     int    `query:"page" doc:"Page number"`
	// 	PageSize int    `query:"page_size" doc:"Number of items per page"`
	// 	Search   string `query:"search" doc:"Search query"`
	// }) (*dtos.ListClientResponse, error) {
	// 	query := core_repositories_ports.ClientQuery{
	// 		Page:     &input.Page,
	// 		PageSize: &input.PageSize,
	// 		Search:   &input.Search,
	// 	}
	// 	clients, err := clientService.ListClients(ctx, query)

	// 	if err != nil {
	// 		return nil, server_errors.AppErrorToHumaError(*err)
	// 	}

	// 	resp := &dtos.ListClientResponse{
	// 		Body: clients,
	// 	}
	// 	return resp, nil

	// })

	// huma.Register(api, huma.Operation{
	// 	OperationID:   "client-detail",
	// 	Method:        http.MethodGet,
	// 	Path:          "/client/{id}",
	// 	Description:   "Get a client by ID",
	// 	Tags:          []string{"Core", "Client"},
	// 	DefaultStatus: http.StatusOK,
	// 	Security: []map[string][]string{
	// 		{"TokenAuth": {"api"}},
	// 	},
	// }, func(ctx context.Context, input *struct {
	// 	Id string `path:"id" maxLength:"24" example:"671495c80d8e9fd6b3256340" doc:"Client ID"`
	// }) (*dtos.ClientDetailResponse, error) {
	// 	client, err := clientService.GetClientByID(ctx, input.Id)

	// 	if err != nil {
	// 		return nil, server_errors.AppErrorToHumaError(*err)
	// 	}

	// 	resp := &dtos.ClientDetailResponse{
	// 		Body: *client,
	// 	}
	// 	return resp, nil

	// })

	// huma.Register(api, huma.Operation{
	// 	OperationID:   "client-update-partial",
	// 	Method:        http.MethodPatch,
	// 	Path:          "/client/{id}",
	// 	Description:   "Update partial client by ID",
	// 	Tags:          []string{"Core", "Client"},
	// 	DefaultStatus: http.StatusOK,
	// 	Security: []map[string][]string{
	// 		{"TokenAuth": {"api"}},
	// 	},
	// }, func(ctx context.Context, input *dtos.UpdatePartialClientRequest) (*dtos.ClientDetailResponse, error) {

	// 	appContext := context.Background()
	// 	tokenValue := ctx.Value(ctx_config.TokenContextKey)
	// 	token, _ := tokenValue.(*entities.Token)
	// 	appContext = context.WithValue(appContext, ctx_config.ClientContextKey, token.User)

	// 	inputData := dtos.UpdatePartialClientInput{}
	// 	if input.Body.Name.Valid {
	// 		inputData.Name = input.Body.Name
	// 	}

	// 	if input.Body.IsActive != nil {
	// 		inputData.IsActive = input.Body.IsActive
	// 	}

	// 	currentClient, err := clientService.UpdatePartialClient(appContext, input.Id, inputData)

	// 	if err != nil {
	// 		return nil, server_errors.AppErrorToHumaError(*err)
	// 	}

	// 	resp := &dtos.ClientDetailResponse{
	// 		Body: *currentClient,
	// 	}
	// 	return resp, nil

	// })

	return &ret

}
