package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"moonshine/internal/domain"
	"moonshine/internal/repository"
	"moonshine/internal/service"
	"moonshine/internal/util"
	"net/http"
	"os"
	"time"
)

type SignUpForm struct {
	Username string `json:"username" form:"username" valid:"required,length(3|20)"`
	Email    string `json:"email" form:"email" valid:"required,email"`
	Password string `json:"password" form:"password" valid:"required,length(3|20)"`
}

type SignInForm struct {
	Username string `json:"username" form:"username" valid:"required,length(3|20)"`
	Password string `json:"password" form:"password" valid:"required,length(3|20)"`
}

type jwtCustomClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func SignUp(c echo.Context) error {
	form := new(SignUpForm)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if _, err := govalidator.ValidateStruct(form); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := domain.User{
		Username: form.Username,
		Email:    form.Email,
		Password: util.HashPassword(form.Password),
	}

	if _, err := service.CreateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, "Email or username already exists")
	} else {
		return c.JSON(http.StatusOK, "")
	}
}

func SignIn(c echo.Context) error {
	form := new(SignInForm)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if _, err := govalidator.ValidateStruct(form); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userRepo := repository.NewUserRepository()
	user, err := userRepo.FindByUsername(form.Username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "User not found")
	}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Incorrect password")
		} else {
			return generateJwtToken(c, *user)
		}
}

func generateJwtToken(c echo.Context, user domain.User) error {
	t, err := GenerateJwtPayload(user.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func GenerateJwtPayload(id uint) (string, error) {
	claims := jwt.MapClaims{
		"id": float64(id),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return t, err
	}

	return t, nil
}
