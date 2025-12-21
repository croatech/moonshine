package mutations

import (
	"errors"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"moonshine/internal/domain"
	"moonshine/internal/graphql/models"
)

var (
	// Validation errors - returned to user
	errInvalidInput       = errors.New("invalid input")
	errInvalidCredentials = errors.New("invalid credentials")
	errUserAlreadyExists  = errors.New("user already exists")

	// Internal errors - logged, generic message to user
	errInternalError = errors.New("internal server error")
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

func domainUserToGraphQL(user *domain.User) *models.User {
	return &models.User{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Hp:        int(user.Hp),
		Level:     int(user.Level),
		Gold:      int(user.Gold),
		Exp:       int(user.Exp),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func generateJWTToken(id uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"id":  id.String(),
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}
