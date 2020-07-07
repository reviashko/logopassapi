package auth

import (
	"crypto/rand"
	"fmt"
)

//GetNewPassword function
func GetNewPassword() (string, error) {
	guidBytes := make([]byte, 16)
	_, err := rand.Read(guidBytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", guidBytes[0:4]), nil
}
