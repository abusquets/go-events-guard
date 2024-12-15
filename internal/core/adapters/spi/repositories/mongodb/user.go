package mongodb

import (
	"context"
	"eventsguard/internal/app/errors"
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

type userRepository struct {
	userCollection *mongo.Collection
	logger         mylog.Logger
}

func NewUserRepository(client *mongo.Client, config *config.AppConfig) repository_ports.UserRepository {
	collection := client.Database(config.MongoDB).Collection("users")
	logger := mylog.GetLogger()
	return &userRepository{userCollection: collection, logger: logger}
}

func (ur *userRepository) GetByID(ctx context.Context, ID string) (*entities.User, *errors.AppError) {
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		ur.logger.ErrorWithErr("Invalid ID format", err)
		return nil, errors.NewValidationError("Invalid ID format")
	}

	ur.logger.Info(fmt.Sprintf("Querying user with filter: %+v", bson.M{"_id": objectID}))

	var user entities.User
	err = ur.userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError(fmt.Sprintf("User not found with ID: %s", ID))
		}
		ur.logger.ErrorWithErr("Error finding user", err)
		return nil, errors.NewUnexpectedError("Error retrieving user")
	}

	ur.logger.Info(fmt.Sprintf("Found user: %+v", user))
	return &user, nil
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, *errors.AppError) {
	email = strings.ToLower(email)

	var user entities.User
	err := ur.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("User  not found")
		}
		ur.logger.ErrorWithErr("Error finding user", err)
		ur.logger.Info(fmt.Sprintf("Query: %+v", bson.M{"email": email}))
		return nil, errors.NewUnexpectedError("Error retrieving user")
	}

	ur.logger.Info(fmt.Sprintf("Found user: %+v", user))

	return &user, nil
}

func (ur *userRepository) Create(ctx context.Context, userData dtos.CreateUserInput) (*entities.User, *errors.AppError) {
	newUser, err := entities.NewUser(
		userData.Email, userData.FirstName, userData.LastName, userData.Password, false, userData.IsActive,
	)
	if err != nil {
		return nil, errors.NewUnexpectedError("Error creating user")
	}
	result, err := ur.userCollection.InsertOne(ctx, newUser)
	if err != nil {
		ur.logger.ErrorWithErr("Error creating user", err)
		msg := err.Error()
		if strings.Contains(msg, "duplicate key error collection: eventsguard.users") {
			msg = "User already exists"
		}
		return nil, errors.NewUnexpectedError("Error creating user: " + msg)
	}

	var createdUser entities.User
	err = ur.userCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdUser)
	if err != nil {
		return nil, errors.NewUnexpectedError("Error retrieving created user")
	}

	return &createdUser, nil
}

func (ur *userRepository) countDocuments(ctx context.Context, filter bson.M) (int64, *errors.AppError) {
	total, err := ur.userCollection.CountDocuments(ctx, filter)
	if err != nil {
		ur.logger.ErrorWithErr("Error counting Users", err)
		return 0, errors.NewUnexpectedError("Error counting users")
	}
	return total, nil
}

func (ur *userRepository) List(ctx context.Context, query repository_ports.UserQuery) (*pagination.PaginatedResult[entities.User], *errors.AppError) {
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
				{"first_name": bson.M{"$regex": *query.Search, "$options": "i"}},
				{"last_name": bson.M{"$regex": *query.Search, "$options": "i"}},
				{"email": bson.M{"$regex": *query.Search, "$options": "i"}},
			},
		}
	}

	total, error := ur.countDocuments(ctx, filter)
	if error != nil {
		return nil, error
	}

	cursor, err := ur.userCollection.Find(ctx, filter, findOptions)
	if err != nil {
		ur.logger.ErrorWithErr("Error listing Users", err)
		return nil, errors.NewUnexpectedError("Error retrieving users")
	}

	defer cursor.Close(ctx)

	var users []entities.User
	for cursor.Next(ctx) {
		var user entities.User
		if err := cursor.Decode(&user); err != nil {
			ur.logger.ErrorWithErr("Error decoding user", err)
			return nil, errors.NewUnexpectedError("Error decoding user")
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		ur.logger.ErrorWithErr("Cursor error", err)
		return nil, errors.NewUnexpectedError("Error with cursor iteration")
	}

	return &pagination.PaginatedResult[entities.User]{
		Items:      users,
		TotalCount: total,
	}, nil
}

func (ur *userRepository) UpdatePartialUser(
	ctx context.Context,
	ID string,
	userData dtos.UpdatePartialAdminUserInput,
) (*entities.User, *errors.AppError) {

	ur.logger.Debug("Iniciant UpdatePartialUser")

	if !userData.FirstName.Valid && !userData.LastName.Valid &&
		userData.IsActive == nil && userData.IsAdmin == nil && userData.Clients == nil {
		ur.logger.Warn("Error: No data provided for update")
		return nil, errors.NewValidationError("No data provided for update")
	}

	updateFields := bson.M{}
	if userData.FirstName.Valid {
		updateFields["first_name"] = userData.FirstName.String
	}
	if userData.LastName.Valid {
		updateFields["last_name"] = userData.LastName.String
	}
	if userData.IsActive != nil {
		updateFields["is_active"] = *userData.IsActive
	}
	if userData.IsAdmin != nil {
		updateFields["is_admin"] = *userData.IsAdmin
	}

	if userData.Clients != nil {
		updateFields["clients"] = userData.Clients
	}

	if len(updateFields) == 0 {
		ur.logger.Warn("Error: No valid fields provided for update")
		return nil, errors.NewValidationError("No valid fields provided for update")
	}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		ur.logger.ErrorWithErr("Error: Invalid ID format", err)
		return nil, errors.NewValidationError("Invalid ID format")
	}

	ur.logger.Debug("Executant UpdateOne a MongoDB")
	_, err = ur.userCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updateFields})
	if err != nil {
		ur.logger.ErrorWithErr("Error updating user", err)
		return nil, errors.NewUnexpectedError("Error updating user")
	}

	user, appErr := ur.GetByID(ctx, ID)
	if appErr != nil {
		ur.logger.ErrorWithErr("Error recuperant usuari:", appErr)
		return nil, appErr
	}
	ur.logger.Info(fmt.Sprintf("Usuari actualitzat: %+v", user))
	return user, nil
}

// func (ur *userRepository) validateClientReferences(
// 	ctx context.Context,
// 	clientIDs []utils_entities.ID,
// ) *errors.AppError {
// 	for _, clientID := range clientIDs {

// 	//     objectID, err := primitive.ObjectIDFromHex(clientID)
// 	//     if err != nil {
// 	//         return errors.NewValidationError("Invalid client ID format")
// 	//     }

// 	//     var client entities.Client
// 	//     err = clientCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&client)
// 	//     if err != nil {
// 	//         if err == mongo.ErrNoDocuments {
// 	//             return errors.NewValidationError(fmt.Sprintf("Client with ID %s does not exist", clientID))
// 	//         }
// 	//         return errors.NewUnexpectedError("Error verifying client existence")
// 	//     }
// 	}
// 	return nil
// }
