package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	byteArr := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(byteArr, 12)
	if err != nil {
		return "", fmt.Errorf("hashing error: %v", err)
	}

	return string(hash), nil
}

func checkPassword(password string, hashed string) bool {
	currentPassword, hashedPassword := []byte(password), []byte(hashed)
	if err := bcrypt.CompareHashAndPassword(hashedPassword, currentPassword); err != nil {
		return false
	}
	return true
}
