package handlers

import (
	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"moonshine/models"
	"moonshine/modules/database"
	"moonshine/modules/support"
	services "moonshine/services/users"
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

	user := models.User{
		Username: form.Username,
		Email:    form.Email,
		Password: support.HashPassword(form.Password),
	}

	if _, err := services.CreateUser(&user); err != nil {
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

	user := models.User{}
	if database.Connection().Where("username = ?", form.Username).First(&user).RecordNotFound() {
		return c.JSON(http.StatusUnauthorized, "User not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Incorrect password")
	} else {
		return generateJwtToken(c, user)
	}
}

func generateJwtToken(c echo.Context, user models.User) error {
	t, err := GenerateJwtPayload(user.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func GenerateJwtPayload(id uint) (string, error) {
	claims := &jwtCustomClaims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return t, err
	}

	return t, nil
}
