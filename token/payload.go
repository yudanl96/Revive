package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// payload data for the token
type Payload struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	IssuedTime  time.Time `json:"issued_time"`
	ExpiredTime time.Time `json:"expired_time"`
}

var (
	ErrInvalidToken = errors.New("token invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:          tokenID,
		Username:    username,
		IssuedTime:  time.Now(),
		ExpiredTime: time.Now().Add(duration),
	}

	return payload, nil

}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredTime) {
		return ErrExpiredToken
	}
	return nil
}
