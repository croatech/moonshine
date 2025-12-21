package mutations

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"moonshine/internal/graphql/models"
	"moonshine/internal/repository"
)

// SignIn authenticates an existing user
func SignIn(ctx context.Context, userRepo *repository.UserRepository, input models.SignInInput) (*models.AuthPayload, error) {
	// Validation errors are returned directly to user
	if err := validateSignInInput(input); err != nil {
		return nil, err
	}

	user, err := userRepo.FindByUsername(input.Username)
	if err != nil {
		// User not found - return invalid credentials (don't reveal if user exists)
		if errors.Is(err, repository.ErrUserNotFound) || errors.Is(err, sql.ErrNoRows) {
			return nil, errInvalidCredentials
		}
		// Database errors are internal - log and return generic error
		log.Printf("ERROR: Failed to find user by username: %v", err)
		return nil, errInternalError
	}

	// Check that password was read from database
	if len(user.Password) == 0 {
		log.Printf("ERROR: Password not found for user: %s", user.Username)
		return nil, errInternalError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		// Invalid password - return invalid credentials (don't reveal which field is wrong)
		return nil, errInvalidCredentials
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		log.Printf("ERROR: Failed to generate JWT token: %v", err)
		return nil, errInternalError
	}

	return &models.AuthPayload{
		Token: token,
		User:  domainUserToGraphQL(user),
	}, nil
}
