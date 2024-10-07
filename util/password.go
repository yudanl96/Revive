package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func MatchPassword(hashedPW, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
