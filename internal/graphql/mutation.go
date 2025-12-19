package graphql

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"moonshine/internal/domain"
	"moonshine/internal/graphql/models"
	"moonshine/internal/repository"
	"moonshine/internal/util"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) SignUp(ctx context.Context, input models.SignUpInput) (*models.AuthPayload, error) {
	if err := validateSignUpInput(input); err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		return nil, errPasswordProcessing
	}

	user := &domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := r.userRepo.Create(user); err != nil {
		return nil, repository.ErrUserExists
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		return nil, errTokenGeneration
	}

	return &models.AuthPayload{
		Token: token,
		User:  domainUserToGraphQL(user),
	}, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, input models.SignInInput) (*models.AuthPayload, error) {
	if err := validateSignInInput(input); err != nil {
		return nil, err
	}

	user, err := r.userRepo.FindByUsername(input.Username)
	if err != nil {
		return nil, errInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errInvalidCredentials
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		return nil, errTokenGeneration
	}

	return &models.AuthPayload{
		Token: token,
		User:  domainUserToGraphQL(user),
	}, nil
}
