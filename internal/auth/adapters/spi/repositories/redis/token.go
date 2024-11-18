package redis

import (
	"encoding/json"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/dtos"
	"eventsguard/internal/infrastructure/mylog"
	"fmt"

	"context"

	"eventsguard/internal/auth/constants"
	"eventsguard/internal/auth/domain/entities"
	port_repositories "eventsguard/internal/auth/domain/ports/repositories"

	"github.com/go-redis/redis/v8"
)

type redisTokenRepository struct {
	redisClient *redis.Client
	ctx         context.Context
	logger      mylog.Logger
}

func NewRedisTokenRepository(redisClient *redis.Client) port_repositories.TokenRepository {
	ctx := context.Background()
	return &redisTokenRepository{
		redisClient: redisClient,
		ctx:         ctx,
		logger:      mylog.GetLogger(),
	}
}

func (r *redisTokenRepository) prefix(key string) string {
	return fmt.Sprintf("fc_token:%s", key)
}

func (r *redisTokenRepository) Save(token dtos.CreateRawTokenDTO) *errors.AppError {
	tokenEntry := r.prefix(fmt.Sprintf("token:%s", token.Token))
	userEntry := r.prefix(fmt.Sprintf("user:%s", token.UserID))

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		r.logger.ErrorWithErr("Error trying to marshal token", err)
		return errors.NewUnexpectedError(err.Error())
	}

	err = r.redisClient.Set(r.ctx, tokenEntry, tokenJSON, 0).Err()
	if err != nil {
		r.logger.ErrorWithErr("Error saving token entry in redis", err)
		return errors.NewUnexpectedError(err.Error())
	}

	err = r.redisClient.HSet(r.ctx, userEntry, string(token.Device), tokenEntry).Err()
	if err != nil {
		r.logger.ErrorWithErr("Error saving user entry in redis", err)
		return errors.NewUnexpectedError(err.Error())
	}
	return nil
}

func (r *redisTokenRepository) findTokenByRedisKey(key string) (*entities.RawToken, *errors.AppError) {
	tokenData, _ := r.redisClient.Get(r.ctx, key).Result()
	// if err != nil {
	//  return nil, errors.NewUnexpectedError(err.Error())
	// }
	if tokenData == "" {
		return nil, nil
	}
	r.logger.Debug("Raw Token data from redis")
	r.logger.With("token_data", tokenData)

	var rawToken entities.RawToken
	err := json.Unmarshal([]byte(tokenData), &rawToken)
	if err != nil {
		return nil, errors.NewUnexpectedError(err.Error())
	}

	return &rawToken, nil
}

func (r *redisTokenRepository) FindByToken(tokenKey string) (*entities.RawToken, *errors.AppError) {
	return r.findTokenByRedisKey(r.prefix(fmt.Sprintf("token:%s", tokenKey)))
}

func (r *redisTokenRepository) FindByUserID(userID string) (map[constants.TokenDevice]*entities.RawToken, *errors.AppError) {
	tokens := make(map[constants.TokenDevice]*entities.RawToken)
	key := r.prefix(fmt.Sprintf("user:%s", userID))
	deviceTokens, err := r.redisClient.HGetAll(r.ctx, key).Result()
	if err != nil {
		return nil, errors.NewUnexpectedError(err.Error())
	}

	for device, tokenKey := range deviceTokens {
		if !constants.IsValidTokenDevice(device) {
			r.logger.Error(fmt.Sprintf("Invalid device value: %s", device))
			continue
		}

		tokenData, err := r.redisClient.Get(r.ctx, tokenKey).Result()
		if err != nil {
			r.logger.ErrorWithErr("Error trying to get token data", err)
			return nil, errors.NewUnexpectedError(err.Error())
		}

		var rawToken entities.RawToken
		err = json.Unmarshal([]byte(tokenData), &rawToken)
		if err != nil {
			r.logger.ErrorWithErr("Error trying to unmarshal token data", err)
			return nil, errors.NewUnexpectedError(err.Error())
		}

		tokens[constants.TokenDevice(device)] = &rawToken
	}

	return tokens, nil
}

func (r *redisTokenRepository) FindByUserIDAndDevice(userID string, device constants.TokenDevice) (*entities.RawToken, *errors.AppError) {
	deviceTokens, err := r.redisClient.HGetAll(r.ctx, r.prefix(fmt.Sprintf("user:%s", userID))).Result()
	if err != nil {
		r.logger.ErrorWithErr("Error trying to get user tokens", err)
		return nil, errors.NewUnexpectedError(err.Error())
	}

	if redisKey, ok := deviceTokens[string(device)]; ok {
		return r.findTokenByRedisKey(redisKey)
	}

	return nil, nil
}

func (r *redisTokenRepository) DeleteByToken(tokenKey string) (bool, *errors.AppError) {
	tokenData, err := r.redisClient.Get(r.ctx, r.prefix(fmt.Sprintf("token:%s", tokenKey))).Result()
	if err != nil {
		r.logger.ErrorWithErr("Error trying to get token data", err)
		return false, errors.NewUnexpectedError(err.Error())
	}

	var rawToken entities.RawToken
	err = json.Unmarshal([]byte(tokenData), &rawToken)
	if err != nil {
		r.logger.ErrorWithErr("Error trying to unmarshal token data", err)
		return false, errors.NewUnexpectedError(err.Error())
	}

	err = r.redisClient.Del(r.ctx, r.prefix(fmt.Sprintf("token:%s", tokenKey))).Err()
	if err != nil {
		r.logger.ErrorWithErr("Error trying to delete token entry", err)
		return false, errors.NewUnexpectedError(err.Error())
	}

	err = r.redisClient.HDel(r.ctx, r.prefix(fmt.Sprintf("user:%s", rawToken.UserID)), tokenKey).Err()
	if err != nil {
		return false, errors.NewUnexpectedError(err.Error())
	}

	return true, nil
}

func (r *redisTokenRepository) DeleteByUserID(userID string) *errors.AppError {
	deviceTokens, err := r.redisClient.HGetAll(r.ctx, r.prefix(fmt.Sprintf("user:%s", userID))).Result()
	if err != nil {
		r.logger.ErrorWithErr("Error trying to get user tokens", err)
		return errors.NewUnexpectedError(err.Error())
	}

	for _, tokenKey := range deviceTokens {
		err := r.redisClient.Del(r.ctx, tokenKey).Err()
		if err != nil {
			r.logger.ErrorWithErr("Error trying to delete token entry", err)
			return errors.NewUnexpectedError(err.Error())
		}
	}

	err = r.redisClient.Del(r.ctx, r.prefix(fmt.Sprintf("user:%s", userID))).Err()
	if err != nil {
		r.logger.ErrorWithErr("Error trying to delete user entry", err)
		return errors.NewUnexpectedError(err.Error())
	}
	return nil

}
