package val

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile("^[a-zA-Z0-9_]+$").MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string,minLenght int,maxLength int) error {
	n := len(value)
	if n < minLenght || n > maxLength {
		return errors.New("string length must be between minLenght and maxLength")
	}
	return nil
}

// ValidateUsername validates the username.
func ValidateUsername(value string) error {
	err:= ValidateString(value, 6, 100)
	if err != nil {
		return err
	}
	if !isValidUsername(value) {
		return errors.New("username is not valid")
	} 
	return nil
}
// ValidateFullName validates the full name.
func ValidateFullName(value string) error {
	err:= ValidateString(value, 3, 100)
	if err != nil {
		return err
	}
	if !isValidFullName(value) {
		return errors.New("full name is not valid")
	} 
	return nil
}
// ValidateEmail validates the email.
func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _,err:=mail.ParseAddress(value); err!= nil {
		return errors.New("email is not valid")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func ValidateSecretCode(value string) error {
	return ValidateString(value, 32, 128)
}
