package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewID(t *testing.T) {
	// Test that NewID generates a valid ObjectID in hexadecimal format
	id := NewID()
	_, err := primitive.ObjectIDFromHex(id.String())
	assert.NoError(t, err, "NewID should return a valid ObjectID in hex format")
}

func TestMarshalBSONValue(t *testing.T) {
	// Generate a new ID
	id := NewID()
	objectID, _ := primitive.ObjectIDFromHex(id.String())

	// Marshal the ID
	bsonType, data, err := id.MarshalBSONValue()

	// Check that the type is correct
	assert.NoError(t, err, "MarshalBSONValue should not return an error")
	assert.Equal(t, bson.TypeObjectID, bsonType, "BSON type should be ObjectID")
	assert.Equal(t, objectID[:], data, "BSON data should match the ObjectID bytes")
}

func TestUnmarshalBSONValue(t *testing.T) {
	// Generate a new ObjectID
	objectID := primitive.NewObjectID()
	id := ID("")

	// Prepare data and type for UnmarshalBSONValue
	bsonType := bson.TypeObjectID
	data := objectID[:]

	// Unmarshal the ID
	err := id.UnmarshalBSONValue(bsonType, data)

	// Check results
	assert.NoError(t, err, "UnmarshalBSONValue should not return an error")
	assert.Equal(t, objectID.Hex(), id.String(), "ID should match the original ObjectID hex value")
}

func TestUnmarshalBSONValueInvalidType(t *testing.T) {
	// Generate a new ObjectID
	objectID := primitive.NewObjectID()
	id := ID("")

	// Use an invalid BSON type for testing
	bsonType := bson.TypeString
	data := objectID[:]

	// Unmarshal with an invalid type
	err := id.UnmarshalBSONValue(bsonType, data)

	// Check that an error is returned
	assert.Error(t, err, "UnmarshalBSONValue should return an error for invalid BSON type")
}

func TestString(t *testing.T) {
	// Test that the String method returns the correct value
	id := NewID()
	assert.Equal(t, string(id), id.String(), "String method should return the ID as a string")
}
