package graphql

import (
	"context"
	"errors"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"

	"moonshine/internal/domain"
	"moonshine/internal/graphql/models"
	"moonshine/internal/util"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) SignUp(ctx context.Context, input models.SignUpInput) (*models.AuthPayload, error) {
	if err := validateSignUpInput(input); err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	user := &domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := r.userRepo.Create(user); err != nil {
		return nil, errors.New("email or username already exists")
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
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
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("incorrect password")
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.AuthPayload{
		Token: token,
		User:  domainUserToGraphQL(user),
	}, nil
}

func validateSignUpInput(input models.SignUpInput) error {
	req := struct {
		Username string `valid:"required,length(3|20)"`
		Email    string `valid:"required,email"`
		Password string `valid:"required,length(3|20)"`
	}{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}
	return nil
}

func validateSignInInput(input models.SignInInput) error {
	req := struct {
		Username string `valid:"required,length(3|20)"`
		Password string `valid:"required,length(3|20)"`
	}{
		Username: input.Username,
		Password: input.Password,
	}

	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}
	return nil
}
