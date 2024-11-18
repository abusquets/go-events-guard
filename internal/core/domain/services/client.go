package services

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/core/domain/entities"
	"eventsguard/internal/core/domain/ports/repositories"
	core_repository_ports "eventsguard/internal/core/domain/ports/repositories"
	core_service_ports "eventsguard/internal/core/domain/ports/services"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/infrastructure/config"
	"eventsguard/internal/infrastructure/mylog"
	"eventsguard/internal/infrastructure/signals"
	"eventsguard/internal/utils/dtos/pagination"
)

type clientService struct {
	clientRepository core_repository_ports.ClientRepository
	logger           mylog.Logger
	asignalsBus      signals.SignalsBus
}

func NewClientService(
	cfg *config.AppConfig,
	clientRepository core_repository_ports.ClientRepository,
	asignalsBus signals.SignalsBus,
) core_service_ports.ClientService {
	return clientService{
		clientRepository: clientRepository,
		logger:           mylog.GetLogger(),
		asignalsBus:      asignalsBus,
	}
}

func (us clientService) CreateClient(ctx context.Context, clientData dtos.CreateClientInput) (*entities.Client, *errors.AppError) {
	return us.clientRepository.Create(ctx, clientData)
}

func (us clientService) GetClientByEmail(ctx context.Context, Email string) (*entities.Client, *errors.AppError) {
	return us.clientRepository.GetByEmail(ctx, Email)
}

func (us clientService) GetClientByID(ctx context.Context, ID string) (*entities.Client, *errors.AppError) {
	return us.clientRepository.GetByID(ctx, ID)
}

func (us clientService) ListClients(ctx context.Context, query repositories.ClientQuery) (*pagination.PaginatedResult[entities.Client], *errors.AppError) {
	return us.clientRepository.List(ctx, query)
}

func (us clientService) UpdatePartialClient(
	ctx context.Context,
	ID string,
	clientData dtos.UpdatePartialClientInput,
) (*entities.Client, *errors.AppError) {
	// responsible, ok := ctx.Value(ctx_config.ClientContextKey).(auth_entities.FakeClient)
	// if ok {
	// 	if !responsible.IsAdmin && responsible.ID != ID {
	// 		return nil, errors.NewPermissionDeniedError(
	// 			"Forbidden: client cannot update other clients",
	// 		)
	// 	}
	// }
	client, error := us.clientRepository.UpdatePartialClient(ctx, ID, clientData)
	if error == nil {
		us.asignalsBus.Emit("client:updated", client.ID)
	}
	return client, error
}
