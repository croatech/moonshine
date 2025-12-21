package mutations

import (
	"context"
	"errors"
	"log"

	"moonshine/internal/domain"
	"moonshine/internal/graphql/models"
	"moonshine/internal/repository"
	"moonshine/internal/util"
)

// SignUp creates a new user account
func SignUp(ctx context.Context, userRepo *repository.UserRepository, input models.SignUpInput) (*models.AuthPayload, error) {
	if err := validateSignUpInput(input); err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		log.Printf("ERROR: Failed to hash password: %v", err)
		return nil, errInternalError
	}

	user := &domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
		
	}

	if err := userRepo.Create(user); err != nil {
		// Check if error is ErrUserExists - this is a validation error, return to user
		if errors.Is(err, repository.ErrUserExists) {
			return nil, errUserAlreadyExists
		}
		// All other database errors are internal - log and return generic error
		log.Printf("ERROR: Failed to create user: %v", err)
		return nil, errInternalError
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
