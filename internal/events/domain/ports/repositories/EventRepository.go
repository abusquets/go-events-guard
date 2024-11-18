package repositories

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/events/domain/entities"
	"eventsguard/internal/events/dtos"
	"eventsguard/internal/utils/dtos/pagination"
	utils_entities "eventsguard/internal/utils/entities"
)

type EventQuery struct {
	Page     *int
	PageSize *int
	Search   *string
	Type     *string
}

type EventRepository interface {
	GetByID(ctx context.Context, clientID utils_entities.ID, eventID string) (*entities.Event, *errors.AppError)

	Create(ctx context.Context, clientID utils_entities.ID, eventData dtos.CreateEventInput) (*entities.Event, *errors.AppError)

	List(ctx context.Context, clientID utils_entities.ID, query EventQuery) (*pagination.PaginatedResult[entities.Event], *errors.AppError)
}
