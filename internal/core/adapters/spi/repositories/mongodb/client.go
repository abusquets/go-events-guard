package mongodb

import (
	"context"
	"errors"
	app_errors "eventsguard/internal/app/errors"
	"eventsguard/internal/core/domain/entities"
	repository_ports "eventsguard/internal/core/domain/ports/repositories"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/infrastructure/config"
	"eventsguard/internal/infrastructure/mylog"
	"eventsguard/internal/utils/dtos/pagination"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type clientRepository struct {
	clientCollection *mongo.Collection
	logger           mylog.Logger
}

func NewClientRepository(client *mongo.Client, config *config.AppConfig) repository_ports.ClientRepository {
	collection := client.Database(config.MongoDB).Collection("clients")
	logger := mylog.GetLogger()
	return &clientRepository{clientCollection: collection, logger: logger}
}

func (ur *clientRepository) GetByID(ctx context.Context, ID string) (*entities.Client, *app_errors.AppError) {
	// Convert the ID from string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		ur.logger.ErrorWithErr("Invalid ID format", err)
		return nil, app_errors.NewValidationError("Invalid ID format")
	}

	ur.logger.Info(fmt.Sprintf("Querying client with filter: %+v", bson.M{"_id": objectID}))

	var client entities.Client
	err = ur.clientCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&client)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			msg := fmt.Sprintf("Client not found with ID: %s", ID)
			ur.logger.Info(msg)
			return nil, app_errors.NewNotFoundError(msg)
		}

		if cmdErr, ok := err.(*mongo.CommandError); ok {
			ur.logger.Error(fmt.Sprintf("Mongo command error: Code: %d, Message: %s", cmdErr.Code, cmdErr.Message))
		}

		ur.logger.ErrorWithErr("Error finding client", err)
		return nil, app_errors.NewUnexpectedError("Error retrieving client")
	}

	ur.logger.Info(fmt.Sprintf("Found client: %+v", client))
	return &client, nil
}

func (ur *clientRepository) Create(ctx context.Context, clientData dtos.CreateClientInput) (*entities.Client, *app_errors.AppError) {
	newClient, err := entities.NewClient(
		clientData.Code, clientData.Name, clientData.IsActive,
	)
	if err != nil {
		return nil, app_errors.NewUnexpectedError("Error creating client")
	}
	result, err := ur.clientCollection.InsertOne(ctx, newClient)
	if err != nil {
		ur.logger.ErrorWithErr("Error creating client", err)
		msg := err.Error()
		if strings.Contains(msg, "duplicate key error collection: eventsguard.clients") {
			msg = "Client already exists"
		}
		return nil, app_errors.NewUnexpectedError("Error creating client: " + msg)
	}

	var createdClient entities.Client
	err = ur.clientCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdClient)
	if err != nil {
		ur.logger.ErrorWithErr("Error retrieving created client", err)
		return nil, app_errors.NewUnexpectedError("Error retrieving created client")
	}

	return &createdClient, nil
}

func (ur *clientRepository) countDocuments(ctx context.Context, filter bson.M) (int64, *app_errors.AppError) {
	total, err := ur.clientCollection.CountDocuments(ctx, filter)
	if err != nil {
		ur.logger.ErrorWithErr("Error counting Clients", err)
		return 0, app_errors.NewUnexpectedError("Error counting clients")
	}
	return total, nil
}

func (ur *clientRepository) List(ctx context.Context, query repository_ports.ClientQuery) (*pagination.PaginatedResult[entities.Client], *app_errors.AppError) {
	findOptions := options.Find()
	if query.Page != nil && *query.Page > 0 && query.PageSize != nil && *query.PageSize > 0 {
		limit := int64(*query.PageSize)
		skip := int64((*query.Page - 1) * *query.PageSize)

		findOptions.SetLimit(limit)
		findOptions.SetSkip(skip)
	} else if query.PageSize != nil && *query.PageSize > 0 {
		findOptions.SetLimit(int64(*query.PageSize))
	}
	filter := bson.M{}
	if query.Search != nil && *query.Search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"code": bson.M{"$regex": *query.Search, "$options": "i"}},
				{"name": bson.M{"$regex": *query.Search, "$options": "i"}},
			},
		}
	}

	total, error := ur.countDocuments(ctx, filter)
	if error != nil {
		return nil, error
	}

	cursor, err := ur.clientCollection.Find(ctx, filter, findOptions)
	if err != nil {
		ur.logger.ErrorWithErr("Error listing Clients", err)
		return nil, app_errors.NewUnexpectedError("Error retrieving clients")
	}

	defer cursor.Close(ctx)

	var clients []entities.Client
	for cursor.Next(ctx) {
		var client entities.Client
		if err := cursor.Decode(&client); err != nil {
			ur.logger.ErrorWithErr("Error decoding client", err)
			return nil, app_errors.NewUnexpectedError("Error decoding client")
		}
		clients = append(clients, client)
	}

	if err := cursor.Err(); err != nil {
		ur.logger.ErrorWithErr("Cursor error", err)
		return nil, app_errors.NewUnexpectedError("Error with cursor iteration")
	}

	return &pagination.PaginatedResult[entities.Client]{
		Items:      clients,
		TotalCount: total,
	}, nil
}

func (ur *clientRepository) UpdatePartialClient(
	ctx context.Context,
	ID string,
	clientData dtos.UpdatePartialClientInput,
) (*entities.Client, *app_errors.AppError) {

	ur.logger.Debug("Iniciant UpdatePartialClient")

	if !clientData.Name.Valid &&
		clientData.IsActive == nil {
		ur.logger.Warn("Error: No data provided for update")
		return nil, app_errors.NewValidationError("No data provided for update")
	}

	// Construir els camps d'update
	updateFields := bson.M{}
	if clientData.Name.Valid {
		updateFields["name"] = clientData.Name.String
	}

	if clientData.IsActive != nil {
		updateFields["is_active"] = *clientData.IsActive
	}

	if len(updateFields) == 0 {
		ur.logger.Warn("Error: No valid fields provided for update")
		return nil, app_errors.NewValidationError("No valid fields provided for update")
	}

	ur.logger.Debug("Convertint ID a ObjectID")
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		ur.logger.ErrorWithErr("Error: Invalid ID format", err)
		return nil, app_errors.NewValidationError("Invalid ID format")
	}

	ur.logger.Debug("Executant UpdateOne a MongoDB")
	_, err = ur.clientCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updateFields})
	if err != nil {
		ur.logger.ErrorWithErr("Error updating client", err)
		return nil, app_errors.NewUnexpectedError("Error updating client")
	}

	client, appErr := ur.GetByID(ctx, ID)
	if appErr != nil {
		ur.logger.ErrorWithErr("Error recuperant client:", appErr)
		return nil, appErr
	}
	ur.logger.Info(fmt.Sprintf("Client actualitzat: %+v", client))
	return client, nil
}
