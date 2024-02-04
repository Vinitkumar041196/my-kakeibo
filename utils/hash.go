package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GetHash(str string) (string, error) {
	// Generate a password hash
	hashedStr, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating password hash:", err)
		return "", err
	}
	return string(hashedStr), nil
}

func CompareStringWithHash(str, hashedStr string) error {
	// Compare a plain-text password with the hash
	return bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
}
