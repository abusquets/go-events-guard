package usecases

import (
	"context"
	"eventsguard/internal/app/errors"
	usecases_ports "eventsguard/internal/inbound_events/domain/ports/usecases"
	"eventsguard/internal/inbound_events/dtos"
	"eventsguard/internal/infrastructure/mylog"
)

type createInboundEvent struct {
	logger mylog.Logger
}

func NewCreateInboundEvent() usecases_ports.CreateInboundEventUseCase {
	return createInboundEvent{
		logger: mylog.GetLogger(),
	}
}

func (u createInboundEvent) Execute(
	ctx context.Context,
	data dtos.CreateInboundEventInput,
) (error *errors.AppError) {

	return nil
}
