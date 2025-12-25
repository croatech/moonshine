package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/dto"
	"moonshine/internal/api/services"
	"moonshine/internal/repository"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(db *sqlx.DB) *AuthHandler {
	userRepo := repository.NewUserRepository(db)
	avatarRepo := repository.NewAvatarRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	authService := services.NewAuthService(userRepo, avatarRepo, locationRepo)

	return &AuthHandler{
		authService: authService,
	}
}

type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string    `json:"token"`
	User  *dto.User `json:"user"`
}

// SignUp handles user registration
// POST /api/auth/signup
func (h *AuthHandler) SignUp(c echo.Context) error {
	var req SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	serviceInput := services.SignUpInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := h.authService.SignUp(c.Request().Context(), serviceInput)
	if err != nil {
		if err == services.ErrUserAlreadyExists {
			return c.JSON(http.StatusConflict, map[string]string{"error": "user already exists"})
		}
		if err == services.ErrInvalidInput {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  dto.UserFromDomain(user),
	})
}

// SignIn handles user authentication
// POST /api/auth/signin
func (h *AuthHandler) SignIn(c echo.Context) error {
	var req SignInRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	serviceInput := services.SignInInput{
		Username: req.Username,
		Password: req.Password,
	}

	user, token, err := h.authService.SignIn(c.Request().Context(), serviceInput)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}
		if err == services.ErrInvalidInput {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  dto.UserFromDomain(user),
	})
}
