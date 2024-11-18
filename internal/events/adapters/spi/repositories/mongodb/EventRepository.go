package mongodb

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/events/domain/dtos"
	"eventsguard/internal/events/domain/entities"
	repository_ports "eventsguard/internal/events/domain/ports/repositories"
	"eventsguard/internal/infrastructure/config"
	"eventsguard/internal/infrastructure/mylog"
	"eventsguard/internal/utils/dtos/pagination"
	utils_entities "eventsguard/internal/utils/entities" // Import del tipus ID
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventRepository struct {
	eventCollection *mongo.Collection
	logger          mylog.Logger
}

func NewEventRepository(client *mongo.Client, config *config.AppConfig) repository_ports.EventRepository {
	collection := client.Database(config.MongoDB).Collection("events")
	logger := mylog.GetLogger()
	return &eventRepository{eventCollection: collection, logger: logger}
}

// GetByID retorna un `Event` per ID i `ClientID`
func (er *eventRepository) GetByID(ctx context.Context, clientID utils_entities.ID, eventID string) (*entities.Event, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		er.logger.ErrorWithErr("Invalid ID format", err)
		return nil, errors.NewValidationError("Invalid ID format")
	}

	// Filtrar per ClientID i EventID
	filter := bson.M{"_id": objectID, "client_id": clientID}

	er.logger.Info(fmt.Sprintf("Querying event with filter: %+v", filter))

	var event entities.Event
	err = er.eventCollection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError(fmt.Sprintf("Event not found with ID: %s", eventID))
		}
		er.logger.ErrorWithErr("Error finding event", err)
		return nil, errors.NewUnexpectedError("Error retrieving event")
	}

	return &event, nil
}

// Create crea un nou `Event` associat a un `ClientID`
func (er *eventRepository) Create(ctx context.Context, clientID utils_entities.ID, eventData dtos.CreateEventInput) (*entities.Event, *errors.AppError) {
	newEvent, err := entities.NewEvent(
		eventData.Type, clientID, eventData.Version, eventData.Payload, eventData.SendAt,
	)
	if err != nil {
		return nil, errors.NewUnexpectedError("Error creating event")
	}

	newEvent.ClientID = clientID

	result, err := er.eventCollection.InsertOne(ctx, newEvent)
	if err != nil {
		er.logger.ErrorWithErr("Error creating event", err)
		msg := err.Error()
		if strings.Contains(msg, "duplicate key error collection") {
			msg = "Event already exists"
		}
		return nil, errors.NewUnexpectedError("Error creating event: " + msg)
	}

	var createdEvent entities.Event
	err = er.eventCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdEvent)
	if err != nil {
		return nil, errors.NewUnexpectedError("Error retrieving created event")
	}

	return &createdEvent, nil
}

func (er *eventRepository) countDocuments(ctx context.Context, filter bson.M) (int64, *errors.AppError) {
	total, err := er.eventCollection.CountDocuments(ctx, filter)
	if err != nil {
		er.logger.ErrorWithErr("Error counting Clients", err)
		return 0, errors.NewUnexpectedError("Error counting clients")
	}
	return total, nil
}

// List retorna una llista paginada d'events filtrats per `ClientID`
func (er *eventRepository) List(ctx context.Context, clientID utils_entities.ID, query repository_ports.EventQuery) (*pagination.PaginatedResult[entities.Event], *errors.AppError) {
	findOptions := options.Find()
	if query.Page != nil && *query.Page > 0 && query.PageSize != nil && *query.PageSize > 0 {
		limit := int64(*query.PageSize)
		skip := int64((*query.Page - 1) * *query.PageSize)
		findOptions.SetLimit(limit)
		findOptions.SetSkip(skip)
	}

	// Filtrar sempre per `ClientID`
	filter := bson.M{"client_id": clientID}
	if query.Search != nil && *query.Search != "" {
		filter["$or"] = []bson.M{
			{"type": bson.M{"$regex": *query.Search, "$options": "i"}},
			{"payload": bson.M{"$regex": *query.Search, "$options": "i"}},
		}
	}

	total, countErr := er.countDocuments(ctx, filter)
	if countErr != nil {
		return nil, countErr
	}

	cursor, err := er.eventCollection.Find(ctx, filter, findOptions)
	if err != nil {
		er.logger.ErrorWithErr("Error listing events", err)
		return nil, errors.NewUnexpectedError("Error retrieving events")
	}
	defer cursor.Close(ctx)

	var events []entities.Event
	for cursor.Next(ctx) {
		var event entities.Event
		if err := cursor.Decode(&event); err != nil {
			er.logger.ErrorWithErr("Error decoding event", err)
			return nil, errors.NewUnexpectedError("Error decoding event")
		}
		events = append(events, event)
	}

	return &pagination.PaginatedResult[entities.Event]{
		Items:      events,
		TotalCount: total,
	}, nil
}
