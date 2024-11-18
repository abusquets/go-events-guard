package services

import (
	"encoding/json"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/domain/entities"
	"time"

	auth_repositories_ports "eventsguard/internal/auth/domain/ports/repositories"
	"eventsguard/internal/auth/domain/ports/services"
	"eventsguard/internal/infrastructure/config"
	"eventsguard/internal/infrastructure/mylog"

	"github.com/google/uuid"

	"eventsguard/internal/auth/constants"

	"eventsguard/internal/auth/dtos"
)

type tokenService struct {
	repository            auth_repositories_ports.TokenRepository
	expiringTokenDuration time.Duration
	logger                mylog.Logger
}

func NewTokenService(
	repository auth_repositories_ports.TokenRepository,
	config *config.AppConfig,
) services.TokenService {

	return &tokenService{
		repository:            repository,
		expiringTokenDuration: time.Duration(config.TokenExpiringDuration) * time.Second,
		logger:                mylog.GetLogger(),
	}
}

func (s *tokenService) CreateByUser(
	user entities.FakeUser,
	device constants.TokenDevice,
	expirable bool,
	renew *bool,
) (*entities.Token, *errors.AppError) {
	if !constants.IsValidTokenDevice(string(device)) {
		s.logger.Error("device - invalid device")
		return nil, errors.NewValidationError("device - invalid device")
	}

	token, error := s.repository.FindByUserIDAndDevice(user.Username, device)

	if error != nil {
		s.logger.Error("error trying to find token by user id and device" + error.Message)
		return nil, error
	}

	var tokenUUID string
	if token != nil && renew != nil && *renew {
		tokenUUID = token.Token
	} else {
		if token != nil {
			s.DeleteByToken(token.Token)
		}
		tokenUUID = s.generateUUID()
	}

	payload := s.makeUserPayload(user)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		s.logger.ErrorWithErr("error trying to marshal payload", err)
		return nil, errors.NewUnexpectedError(err.Error())
	}

	createdAt := time.Now()

	var expiresAt *time.Time
	var expiracy *int
	if expirable {
		expTime := createdAt.Add(s.expiringTokenDuration)
		expiresAt = &expTime
		exp := int(s.expiringTokenDuration.Seconds())
		expiracy = &exp
	}

	dto := dtos.CreateRawTokenDTO{
		Device:    device,
		Token:     tokenUUID,
		Payload:   string(payloadJSON),
		UserID:    user.Username,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
		Expiracy:  expiracy,
	}

	error = s.repository.Save(dto)
	if error != nil {
		s.logger.Error("error trying to save token")
		return nil, error
	}

	return s.FindByTokenKey(tokenUUID)
}

func (s *tokenService) FindByTokenKey(tokenKey string) (*entities.Token, *errors.AppError) {
	rawToken, error := s.repository.FindByToken(tokenKey)
	if error != nil {
		return nil, error
	}
	if rawToken == nil {
		return nil, nil
	}

	return s.fromRawToken(rawToken), nil
}

func (s *tokenService) DeleteByToken(tokenKey string) (bool, *errors.AppError) {
	return s.repository.DeleteByToken(tokenKey)
}

func (s *tokenService) DeleteByUserID(userID string, devices []constants.TokenDevice) *errors.AppError {
	if devices != nil {
		tokens, error := s.repository.FindByUserID(userID)
		if error != nil {
			return error
		}

		for device, rawToken := range tokens {
			if constants.IsValidTokenDevice(string(device)) {
				s.repository.DeleteByToken(rawToken.Token)
			}
		}
	} else {
		return s.repository.DeleteByUserID(userID)
	}
	return nil
}

func (s *tokenService) RenewToken(token entities.Token) *errors.AppError {
	var expiresAt *time.Time
	if token.Expiracy != nil {
		expTime := time.Now().Add(s.expiringTokenDuration)
		expiresAt = &expTime
	}

	payloadJSON, err := json.Marshal(token.User)
	if err != nil {
		return errors.NewUnexpectedError(err.Error())
	}

	dto := dtos.CreateRawTokenDTO{
		Device:    token.Device,
		Token:     token.Token,
		Payload:   string(payloadJSON),
		UserID:    token.UserID,
		ExpiresAt: expiresAt,
		CreatedAt: token.CreatedAt,
		Expiracy:  token.Expiracy,
	}

	return s.repository.Save(dto)
}

func (s *tokenService) IsTokenExpired(token *entities.Token) bool {
	if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
		return true
	}
	return false
}

func (s *tokenService) FindByUserID(userID string) (map[constants.TokenDevice]*entities.Token, *errors.AppError) {
	rawTokens, err := s.repository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	tokens := make(map[constants.TokenDevice]*entities.Token)
	for device, rawToken := range rawTokens {
		tokens[device] = s.fromRawToken(rawToken)
	}

	return tokens, nil
}

func (s *tokenService) fromRawToken(rawToken *entities.RawToken) *entities.Token {
	var user entities.FakeUser
	json.Unmarshal([]byte(rawToken.Payload), &user)

	return &entities.Token{
		Device:    rawToken.Device,
		Token:     rawToken.Token,
		User:      user,
		UserID:    user.Username,
		ExpiresAt: rawToken.ExpiresAt,
		CreatedAt: rawToken.CreatedAt,
		Expiracy:  rawToken.Expiracy,
	}
}

func (s *tokenService) generateUUID() string {
	return uuid.NewString()
}

func (s *tokenService) makeUserPayload(user entities.FakeUser) map[string]interface{} {
	return map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"is_admin":   user.IsAdmin,
	}
}
