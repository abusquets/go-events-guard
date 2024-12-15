package services

import (
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/constants"
	"eventsguard/internal/auth/domain/entities"
	"eventsguard/internal/infrastructure/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	secretKey []byte
}

func NewJWTService(config *config.AppConfig) *JWTService {
	return &JWTService{
		secretKey: []byte(config.JWTSecretKey), // Afegim el secret de configuració per generar els JWT
	}
}

// Funció per crear un token JWT
func (s *JWTService) CreateToken(user entities.FakeUser, device constants.TokenDevice, expirable bool) (string, time.Time, *errors.AppError) {
	claims := jwt.MapClaims{
		"sub":        user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"is_admin":   user.IsAdmin,
		"device":     device,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Durada de 24 hores
	}

	// Si és un token que expira, ajustem el temps d'expiració
	if expirable {
		claims["exp"] = time.Now().Add(time.Duration(24) * time.Hour).Unix()
	}

	// Creem el token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmem el token
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, errors.NewUnexpectedError(err.Error())
	}

	// Establim el temps d'expiració
	expTime := time.Now().Add(time.Duration(24) * time.Hour)

	return tokenString, expTime, nil
}

// Funció per validar el token JWT
func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, *errors.AppError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifiquem que el mètode de signatura sigui correcte
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewUnexpectedError("unexpected signing method")
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, errors.NewUnexpectedError(err.Error())
	}
	return token, nil
}
