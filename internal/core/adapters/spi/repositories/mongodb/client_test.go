package mongodb

import (
	"context"
	"eventsguard/internal/core/domain/entities"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/infrastructure/config"
	utils_entities "eventsguard/internal/utils/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"github.com/jaswdr/faker/v2"
)

func GetClientRepositoryForTest(t *testing.T, mt *mtest.T) *clientRepository {
	mockConfig := &config.AppConfig{
		MongoDB: "testdb",
	}

	repo := NewClientRepository(mt.Client, mockConfig).(*clientRepository)

	return repo
}

func TestGetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer func() {
		if mt.Client != nil {
			mt.Client.Disconnect(context.Background())
		}
	}()

	mt.Run("success - retrieve client by ID", func(mt *mtest.T) {

		require.NotNil(t, mt.Client, "Mock client should not be nil")

		// Create a new ID
		id := utils_entities.NewID()

		// Create the corresponding primitive.ObjectID value
		objID, err := primitive.ObjectIDFromHex(id.String())
		require.NoError(t, err)

		client := entities.Client{
			ID:       id,
			Code:     "12345",
			Name:     "Test Client",
			IsActive: true,
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "testdb.clients", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: objID},
				{Key: "code", Value: client.Code},
				{Key: "name", Value: client.Name},
				{Key: "is_active", Value: client.IsActive},
			}),
		)

		repo := GetClientRepositoryForTest(t, mt)

		ctx := context.Background()
		result, appErr := repo.GetByID(ctx, id.String())

		require.Nil(t, appErr)
		require.NotNil(t, result)
		require.Equal(t, client.ID, result.ID)
		assert.Equal(t, client.Code, result.Code)
		assert.Equal(t, client.Name, result.Name)
	})

	mt.Run("error - invalid ID format", func(mt *mtest.T) {

		repo := GetClientRepositoryForTest(t, mt)

		ctx := context.Background()
		result, appErr := repo.GetByID(ctx, "invalid_id")

		require.NotNil(t, appErr)
		require.Nil(t, result)
		assert.Equal(t, "Invalid ID format", appErr.Message)
	})

	mt.Run("error - ID doesn't exist", func(mt *mtest.T) {
		repo := GetClientRepositoryForTest(t, mt)
		nonExistentID := "675854306484f80032b05567"

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "testdb.clients", mtest.FirstBatch), // 0 documents
		)

		ctx := context.Background()
		result, appErr := repo.GetByID(ctx, nonExistentID)

		assert.Nil(t, result)
		assert.NotNil(t, appErr)
		assert.Equal(t, "Client not found with ID: "+nonExistentID, appErr.Message)
	})

}

func TestCreate(t *testing.T) {
	fake := faker.New()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer func() {
		if mt.Client != nil {
			mt.Client.Disconnect(context.Background())
		}
	}()

	mt.Run("success - create client", func(mt *mtest.T) {
		repo := GetClientRepositoryForTest(t, mt)

		clientData := dtos.CreateClientInput{
			Code:     fake.UUID().V4(),
			Name:     fake.Company().Name(),
			IsActive: true,
		}

		newClient, err := entities.NewClient(clientData.Code, clientData.Name, clientData.IsActive)
		require.NoError(t, err)

		mt.AddMockResponses(
			mtest.CreateSuccessResponse(),
			mtest.CreateCursorResponse(1, "testdb.clients", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: newClient.ID},
				{Key: "code", Value: newClient.Code},
				{Key: "name", Value: newClient.Name},
				{Key: "is_active", Value: newClient.IsActive},
			}),
		)

		ctx := context.Background()
		result, appErr := repo.Create(ctx, clientData)

		require.Nil(t, appErr)
		require.NotNil(t, result)
		assert.Equal(t, clientData.Code, result.Code)
		assert.Equal(t, clientData.Name, result.Name)
		assert.Equal(t, clientData.IsActive, result.IsActive)
	})

	mt.Run("error - duplicate client", func(mt *mtest.T) {
		repo := GetClientRepositoryForTest(t, mt)

		clientData := dtos.CreateClientInput{
			Code:     fake.UUID().V4(),
			Name:     fake.Company().Name(),
			IsActive: true,
		}

		mt.AddMockResponses(
			mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Code:    11000,
				Message: "duplicate key error collection: eventsguard.clients",
			}),
		)

		ctx := context.Background()
		result, appErr := repo.Create(ctx, clientData)

		require.NotNil(t, appErr)
		require.Nil(t, result)
		assert.Equal(t, "Error creating client: Client already exists", appErr.Message)
	})

	mt.Run("error - unexpected error", func(mt *mtest.T) {
		repo := GetClientRepositoryForTest(t, mt)

		clientData := dtos.CreateClientInput{
			Code:     fake.UUID().V4(),
			Name:     fake.Company().Name(),
			IsActive: true,
		}

		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    12345,
				Message: "unexpected error",
			}),
		)

		ctx := context.Background()
		result, appErr := repo.Create(ctx, clientData)

		require.NotNil(t, appErr)
		require.Nil(t, result)
		assert.Equal(t, "Error creating client: unexpected error", appErr.Message)
	})
}
