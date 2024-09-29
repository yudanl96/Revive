package util

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Cannot hash password:", err)
	}
	return string(hashedPassword)
}

func MatchPassword(hashedPW, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
