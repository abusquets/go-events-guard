package usecases

import (
	"context"

	"eventsguard/internal/inbound_events/dtos"

	"eventsguard/internal/app/errors"
)

type CreateInboundEventUseCase interface {
	Execute(ctx context.Context, data dtos.CreateInboundEventInput) *errors.AppError
}
