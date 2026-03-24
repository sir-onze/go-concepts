package main

import (
	"errors"
	"fmt"
)

// --- given ---

var ErrNotFound = errors.New("user not found")

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

func findUser(userID string) (string, error) {
	if userID == "unknown" {
		return "", ErrNotFound
	}
	return "Alice", nil
}

// --- your job ---

func getUser(userID string) (string, error) {
	// TODO:
	// 1. if userID is empty, return a ValidationError{Field: "userID", Message: "cannot be empty"}
	if userID == "" {
		return "", &ValidationError{Field: "userID", Message: "cannot be empty"}
	}
	// 2. call findUser(userID), wrap any error with fmt.Errorf + %w
	user, err := findUser(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	// 3. return the user on success
	return user, nil
}

func main() {
	// test 1 — empty ID
	_, err := getUser("")
	// TODO: use errors.As to check for ValidationError, print the field name
	if err != nil {
		var ve *ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("validation error on field: %s\n", ve.Field)
		} else {
			fmt.Printf("unexpected error: %v\n", err)
		}
	}
	// test 2 — unknown user
	_, err = getUser("unknown")
	// TODO: use errors.Is to check for ErrNotFound, print "user not found"
	if errors.Is(err, ErrNotFound) {
		fmt.Println("user not found")
	} else {
		fmt.Printf("unexpected error: %v\n", err)
	}
	// test 3 — happy path
	user, err := getUser("alice-123")
	// TODO: print the user name
	fmt.Printf("user: %s\n", user)
}
