package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
	"moonshine/internal/util"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInternalError      = errors.New("internal server error")
)

type SignUpInput struct {
	Username string
	Email    string
	Password string
}

type SignInInput struct {
	Username string
	Password string
}

type AuthService struct {
	userRepo     *repository.UserRepository
	avatarRepo   *repository.AvatarRepository
	locationRepo *repository.LocationRepository
}

func NewAuthService(userRepo *repository.UserRepository, avatarRepo *repository.AvatarRepository, locationRepo *repository.LocationRepository) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		avatarRepo:   avatarRepo,
		locationRepo: locationRepo,
	}
}

func (s *AuthService) SignUp(ctx context.Context, input SignUpInput) (*domain.User, string, error) {
	if err := s.validateSignUpInput(input); err != nil {
		return nil, "", err
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		log.Printf("[SignUp ERROR] Failed to hash password - Full error: %+v, Type: %T", err, err)
		return nil, "", ErrInternalError
	}

	// Get default location (moonshine)
	location, err := s.locationRepo.FindStartLocation()
	if err != nil {
		log.Printf("[SignUp ERROR] Failed to find start location - Full error: %+v, Type: %T", err, err)
		return nil, "", ErrInternalError
	}

	// Get first available avatar (optional)
	var avatarID *uuid.UUID
	avatar, err := s.avatarRepo.FindFirst()
	if err == nil && avatar != nil {
		avatarID = &avatar.ID
	}

	user := &domain.User{
		Username:   input.Username,
		Name:       input.Username, // Set Name same as Username
		Email:      input.Email,
		Password:   hashedPassword,
		Attack:     1, // Default attack value (matching DB default)
		Defense:    1, // Default defense value (matching DB default)
		Hp:         20,
		CurrentHp:  20,
		Level:      1,
		Gold:       100,
		Exp:        0,
		FreeStats:  15,
		LocationID: location.ID,
		AvatarID:   avatarID,
	}

	if err := s.userRepo.Create(user); err != nil {
		// Check if error is ErrUserExists - this is a validation error, return to user
		if errors.Is(err, repository.ErrUserExists) {
			log.Printf("[SignUp] User already exists: %s", input.Username)
			return nil, "", ErrUserAlreadyExists
		}
		// All other database errors are internal - log and return generic error
		log.Printf("[SignUp ERROR] Failed to create user - Full error: %+v, Type: %T, Username: %s", err, err, input.Username)
		return nil, "", ErrInternalError
	}

	token, err := s.generateJWTToken(user.ID)
	if err != nil {
		log.Printf("[SignUp ERROR] Failed to generate JWT token - Full error: %+v, Type: %T, UserID: %s", err, err, user.ID)
		return nil, "", ErrInternalError
	}

	return user, token, nil
}

func (s *AuthService) SignIn(ctx context.Context, input SignInInput) (*domain.User, string, error) {
	// Validation errors are returned directly to user
	if err := s.validateSignInInput(input); err != nil {
		return nil, "", err
	}

	user, err := s.userRepo.FindByUsername(input.Username)
	if err != nil {
		// User not found - return invalid credentials (don't reveal if user exists)
		if errors.Is(err, repository.ErrUserNotFound) || errors.Is(err, sql.ErrNoRows) {
			log.Printf("[SignIn] User not found: %s", input.Username)
			return nil, "", ErrInvalidCredentials
		}
		// Database errors are internal - log and return generic error
		log.Printf("[SignIn ERROR] Failed to find user by username - Full error: %+v, Type: %T, Username: %s", err, err, input.Username)
		return nil, "", ErrInternalError
	}

	// Check that password was read from database
	if len(user.Password) == 0 {
		log.Printf("[SignIn ERROR] Password not found for user - Username: %s, UserID: %s", user.Username, user.ID)
		return nil, "", ErrInternalError
	}

	if err := util.CheckPassword(user.Password, input.Password); err != nil {
		// Invalid password - return invalid credentials (don't reveal which field is wrong)
		log.Printf("[SignIn] Invalid password for user: %s", input.Username)
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.generateJWTToken(user.ID)
	if err != nil {
		log.Printf("[SignIn ERROR] Failed to generate JWT token - Full error: %+v, Type: %T, UserID: %s", err, err, user.ID)
		return nil, "", ErrInternalError
	}

	return user, token, nil
}

func (s *AuthService) validateSignUpInput(input SignUpInput) error {
	type signUpValidator struct {
		Username string `valid:"required,length(3|20)"`
		Email    string `valid:"required,email"`
		Password string `valid:"required,length(3|20)"`
	}

	v := signUpValidator{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if _, err := govalidator.ValidateStruct(v); err != nil {
		return ErrInvalidInput
	}
	return nil
}

func (s *AuthService) validateSignInInput(input SignInInput) error {
	type signInValidator struct {
		Username string `valid:"required,length(3|20)"`
		Password string `valid:"required,length(3|20)"`
	}

	v := signInValidator{
		Username: input.Username,
		Password: input.Password,
	}

	if _, err := govalidator.ValidateStruct(v); err != nil {
		return ErrInvalidInput
	}
	return nil
}

func (s *AuthService) generateJWTToken(id uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"id":  id.String(),
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}
