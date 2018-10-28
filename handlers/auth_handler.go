package handlers

import (
	"net/http"
	"sunlight/models"
	"sunlight/modules/database"
	"time"

	"github.com/asaskevich/govalidator"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
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
		Password: models.HashPassword(form.Password),
	}
	if err := database.Connection().Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Email or username already exist")
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
		return c.JSON(http.StatusInternalServerError, "Player not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Incorrect password")
	} else {
		return generateJwtToken(c, user)
	}
}

func generateJwtToken(c echo.Context, user models.User) error {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["expiresAt"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
