package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseErrorMessage(message string) string {
	s := strings.Split(message, "\n")
	var errMessage string
	for _, part := range s {
		// Parse each message and return its parsed form
		step1 := strings.Split(part, ":")[1]  // 'Key' Error
		step2 := strings.Trim(step1, " ")     // 'Key' Error
		step3 := strings.Split(step2, " ")[0] // 'Key'
		errorKey := strings.Trim(step3, "'")  // Key
		msg := errorKey + " cannot be empty;"
		errMessage += msg
	}
	return errMessage
}
