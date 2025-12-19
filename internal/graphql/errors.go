package graphql

import "errors"

var (
	errInvalidCredentials = errors.New("invalid credentials")
	errPasswordProcessing = errors.New("failed to process password")
	errTokenGeneration    = errors.New("failed to generate token")
	errInvalidInput       = errors.New("invalid input")
)


