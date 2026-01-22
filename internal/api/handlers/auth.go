package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/dto"
	"moonshine/internal/api/services"
	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type AuthHandler struct {
	authService   *services.AuthService
	locationRepo  *repository.LocationRepository
}

func NewAuthHandler(db *sqlx.DB) *AuthHandler {
	userRepo := repository.NewUserRepository(db)
	avatarRepo := repository.NewAvatarRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	authService := services.NewAuthService(userRepo, avatarRepo, locationRepo)

	return &AuthHandler{
		authService:  authService,
		locationRepo: locationRepo,
	}
}

type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignInRequest struct {
	Username string `json:"username" validate:"required" example:"admin"`
	Password string `json:"password" validate:"required" example:"password"`
}

type AuthResponse struct {
	Token string    `json:"token"`
	User  *dto.User `json:"user"`
}

// SignUp godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SignUpRequest true "Sign up request"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/auth/signup [post]
func (h *AuthHandler) SignUp(c echo.Context) error {
	var req SignUpRequest
	if err := c.Bind(&req); err != nil {
		return ErrBadRequest(c, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return ErrBadRequest(c, err.Error())
	}

	serviceInput := services.SignUpInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := h.authService.SignUp(c.Request().Context(), serviceInput)
	if err != nil {
		if err == services.ErrUserAlreadyExists {
			return ErrConflict(c, "user already exists")
		}
		if err == services.ErrInvalidInput {
			return ErrBadRequest(c, "invalid input")
		}
		return ErrInternalServerError(c)
	}

	var location *domain.Location
	if user.LocationID != uuid.Nil {
		location, _ = h.locationRepo.FindByID(user.LocationID)
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  dto.UserFromDomain(user, location, false),
	})
}

// SignIn godoc
// @Summary Sign in
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SignInRequest true "Sign in request" example({"username":"admin","password":"password"})
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/signin [post]
func (h *AuthHandler) SignIn(c echo.Context) error {
	var req SignInRequest
	if err := c.Bind(&req); err != nil {
		return ErrBadRequest(c, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return ErrBadRequest(c, err.Error())
	}

	serviceInput := services.SignInInput{
		Username: req.Username,
		Password: req.Password,
	}

	user, token, err := h.authService.SignIn(c.Request().Context(), serviceInput)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			return ErrUnauthorizedWithMessage(c, "invalid credentials")
		}
		if err == services.ErrInvalidInput {
			return ErrBadRequest(c, "invalid input")
		}
		return ErrInternalServerError(c)
	}

	var location *domain.Location
	if user.LocationID != uuid.Nil {
		location, _ = h.locationRepo.FindByID(user.LocationID)
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  dto.UserFromDomain(user, location, false),
	})
}
