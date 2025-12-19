package graphql

import (
	"github.com/asaskevich/govalidator"

	"moonshine/internal/graphql/models"
)

type signUpValidator struct {
	Username string `valid:"required,length(3|20)"`
	Email    string `valid:"required,email"`
	Password string `valid:"required,length(3|20)"`
}

type signInValidator struct {
	Username string `valid:"required,length(3|20)"`
	Password string `valid:"required,length(3|20)"`
}

func validateSignUpInput(input models.SignUpInput) error {
	v := signUpValidator{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if _, err := govalidator.ValidateStruct(v); err != nil {
		return errInvalidInput
	}
	return nil
}

func validateSignInInput(input models.SignInInput) error {
	v := signInValidator{
		Username: input.Username,
		Password: input.Password,
	}

	if _, err := govalidator.ValidateStruct(v); err != nil {
		return errInvalidInput
	}
	return nil
}


