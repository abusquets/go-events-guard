package entities

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ID type that wraps a string
type ID string

// NewID generates a new ID from a primitive.ObjectID
func NewID() ID {
	return ID(primitive.NewObjectID().Hex())
}

// String returns the string representation of the ID
func (id ID) String() string {
	return string(id)
}

func (id ID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	objID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return bson.TypeNull, nil, err // Canviem bsontype.Null per bson.TypeNull
	}

	// Retornem l'ObjectID com a tipus BSON ObjectID
	return bson.TypeObjectID, objID[:], nil
}

// UnmarshalBSON unmarshals the BSON into ID format.
func (id *ID) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	// Check if the type is ObjectID
	if t != bson.TypeObjectID {
		return bson.ErrDecodeToNil
	}

	// Create an ObjectID from the data bytes
	var objID primitive.ObjectID
	copy(objID[:], data)

	// Convert the ObjectID to a hexadecimal string and set the ID
	*id = ID(objID.Hex())
	return nil
}
