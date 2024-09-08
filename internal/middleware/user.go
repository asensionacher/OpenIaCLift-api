package middleware

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash password using the bcrypt hashing algorithm

func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	passwordBytes := []byte(password)
	// Hash password with bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	return string(hashedPasswordBytes), err
}

func checkPassword(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))

	return err == nil
}
