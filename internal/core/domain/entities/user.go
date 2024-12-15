package entities

import (
	"eventsguard/internal/utils/entities"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        entities.ID   `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Email     string        `json:"email" validate:"required,email"`
	FirstName string        `json:"first_name" bson:"first_name" validate:"required"`
	LastName  *string       `json:"last_name" bson:"last_name" validate:"required"`
	Password  string        `json:"-"`
	IsAdmin   bool          `json:"is_admin" bson:"is_admin" validate:"required"`
	IsActive  bool          `json:"is_active" bson:"is_active" validate:"required"`
	Clients   []entities.ID `json:"clients" bson:"clients"`
}

// EncryptPassword encrypts the password using bcrypt
func (u *User) EncryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// VerifyPassword verifies the password against the hashed password
func (u *User) VerifyPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	return err == nil
}

// NewUser creates a new User with a generated ID
func NewUser(email string, firstName string, lastName *string, password string, isAdmin bool, isActive bool) (*User, error) {
	user := &User{
		ID:        entities.NewID(), // Generate a new ID here
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
		IsAdmin:   isAdmin,
		IsActive:  isActive,
	}

	// Encrypt the password before saving
	err := user.EncryptPassword()
	if err != nil {
		return nil, err
	}

	return user, nil
}
