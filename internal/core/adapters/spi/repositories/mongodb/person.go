package mongodb

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/core/domain/entities"
	repository_ports "eventsguard/internal/core/domain/ports/repositories"
	"eventsguard/internal/core/dtos"
	"eventsguard/internal/infrastructure/config"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type personRepository struct {
	personCollection *mongo.Collection
}

func NewPersonRepository(client *mongo.Client, config *config.AppConfig) repository_ports.PersonRepository {
	collection := client.Database(config.MongoDB).Collection("persons")
	return &personRepository{personCollection: collection}
}

// GetByID retrieves a person by their UUID from the database
func (ur *personRepository) GetByID(ctx context.Context, uuid string) (*entities.Person, *errors.AppError) {

	// Convert the UUID string to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, errors.NewValidationError("Invalid UUID format")
	}

	// Define a person variable to hold the result
	var person entities.Person

	// Find the person with the matching ObjectID
	err = ur.personCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&person)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("Person not found")
		}
		return nil, errors.NewUnexpectedError("Error retrieving person")
	}

	return &person, nil
}

// Create inserts a new person into the database
func (ur *personRepository) Create(ctx context.Context, personData dtos.CreatePersonInput) (*entities.Person, *errors.AppError) {

	// Create a new Person entity
	newPerson := entities.Person{
		ID:        primitive.NewObjectID(),
		FirstName: personData.FirstName,
		LastName:  personData.LastName,
	}

	// Insert the person into the database
	result, err := ur.personCollection.InsertOne(ctx, newPerson)
	if err != nil {
		log.Println("Error creating person:", err)
		msg := err.Error()
		if strings.Contains(msg, "duplicate key error collection: eventsguard.persons") {
			msg = "Person already exists"
		}
		return nil, errors.NewUnexpectedError("Error creating person: " + msg)
	}

	var createdPerson entities.Person
	err = ur.personCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdPerson)
	if err != nil {
		return nil, errors.NewUnexpectedError("Error retrieving created person")
	}

	return &createdPerson, nil
}

// List retrieves all persons from the database
func (ur *personRepository) List(ctx context.Context) (*[]entities.Person, *errors.AppError) {

	cursor, err := ur.personCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error finding documents:", err)
		return nil, errors.NewUnexpectedError("Error retrieving persons")
	}
	defer cursor.Close(ctx)

	var persons []entities.Person

	// Itera pels resultats del cursor
	for cursor.Next(ctx) {
		var person entities.Person
		err := cursor.Decode(&person)
		if err != nil {
			log.Println("Error decoding person:", err)
			return nil, errors.NewUnexpectedError("Error decoding person")
		}

		// Afegeix el document desat en la llista d'usuaris
		persons = append(persons, person)
	}

	// Comprova si el cursor ha estat esgotat abans de tancar-lo
	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, errors.NewUnexpectedError("Error with cursor iteration")
	}

	// Retorna la llista d'usuaris (sense fer servir &)
	return &persons, nil
}
