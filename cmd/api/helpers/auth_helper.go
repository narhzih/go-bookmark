package helpers

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash, authOrigin string) (ok bool, err error) {

	if authOrigin == "" || authOrigin == "DEFAULT" {
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
			return false, fmt.Errorf("incorrect password")
		}
		return true, nil
	}
	return false, fmt.Errorf("your account was created using %+v., authentication can only be carried out using the same channel", authOrigin)
}

func ParseErrorMessage(message string) string {
	s := strings.Split(message, "\n")
	var errMessage string

	if message == "EOF" {
		errMessage = "Err: No request body sent;"
	} else if strings.Contains(message, "tag") {
		for _, part := range s {
			// Parse each message and return its parsed form
			step1 := strings.Split(part, ":")[1]  // 'Key' Error
			step2 := strings.Trim(step1, " ")     // 'Key' Error
			step3 := strings.Split(step2, " ")[0] // 'Key'
			errorKey := strings.Trim(step3, "'")  // Key
			msg := errorKey + " cannot be empty;"
			errMessage += msg
		}
	} else {
		// return initial error as it is
		errMessage = message
	}

	return errMessage
}
