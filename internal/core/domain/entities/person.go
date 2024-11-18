package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Person struct represents a person with an ID and its fields
type Person struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"` // MongoDB ID
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  *string            `json:"last_name" bson:"last_name" validate:"required"`
}

// NewPerson creates a new Person with a generated ID
func NewPerson(email, firstName, lastName, password string, isAdmin, isActive bool) (*Person, error) {
	person := &Person{
		ID:        primitive.NewObjectID(), // Generate a new MongoDB ObjectID
		FirstName: firstName,
		LastName:  &lastName,
	}

	return person, nil
}
