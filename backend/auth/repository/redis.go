package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo interface {
	SetRefreshToken(ctx context.Context, userID, token string) error
	GetRefreshToken(ctx context.Context, userID, token string) (bool, error)
	Delete(ctx context.Context, userID string) error
}

type tokenRepo struct {
	rdb             *redis.Client
	refreshTokenTTL time.Duration
	refreshTokenKey string
}

func NewRedisRepo(rdb *redis.Client, refreshTokenTTL time.Duration, refreshTokenKey string) RedisRepo {
	return &tokenRepo{rdb: rdb, refreshTokenTTL: refreshTokenTTL, refreshTokenKey: refreshTokenKey}
}

func (r *tokenRepo) SetRefreshToken(ctx context.Context, userID, token string) error {
	return r.rdb.Set(ctx, r.key(userID), hash(token), r.refreshTokenTTL).Err()
}

func (r *tokenRepo) GetRefreshToken(ctx context.Context, userID, token string) (bool, error) {
	stored, err := r.rdb.Get(ctx, r.key(userID)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stored == hash(token), nil
}

func (r *tokenRepo) Delete(ctx context.Context, userID string) error {
	return r.rdb.Del(ctx, r.key(userID)).Err()
}

func (r *tokenRepo) key(userID string) string {
	return fmt.Sprintf("%s:%s", r.refreshTokenKey, userID)
}

func hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
