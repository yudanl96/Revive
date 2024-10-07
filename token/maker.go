package token

import "time"

// used to manage tokens
type Maker interface {
	// create new token for username for the duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// check if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
