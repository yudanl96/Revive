package redisdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	ID           string
	Username     string
	RefreshToken string
	UserAgent    string
	ClienIp      string
	IsBlocked    bool
	ExpiredTime  time.Time
}

type RedisRepo struct {
	Client *redis.Client
}

func sessionIDKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (r *RedisRepo) CreateSession(ctx context.Context, session Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to encode session: %w", err)
	}

	key := sessionIDKey(session.ID)

	res := r.Client.Set(ctx, key, string(data), time.Until(session.ExpiredTime)) //to be changed later
	if err := res.Err(); err != nil {
		return fmt.Errorf("failed to set: %w", err)
	}

	return nil
}

var ErrNoSession = errors.New("session does not exist")

func (r *RedisRepo) RetrieveSession(ctx context.Context, id string) (Session, error) {
	key := sessionIDKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return Session{}, ErrNoSession
	} else if err != nil {
		return Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	var session Session
	err = json.Unmarshal([]byte(value), &session)
	if err != nil {
		return Session{}, fmt.Errorf("failed to decode session json: %w", err)
	}

	return session, nil
}
