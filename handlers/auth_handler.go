package handlers

import (
	"net/http"
	"sunlight/models"
	"sunlight/modules/database"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type SignUpForm struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type SignInForm struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func SignUp(c echo.Context) error {
	form := new(SignUpForm)

	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user := models.User{
		Username: form.Username,
		Email:    form.Email,
		Password: models.HashPassword(form.Password),
	}

	if err := user.Validate(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := database.Connection().Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Undefined error",
		})
	} else {
		return generateJwtToken(c, user)
	}
}

func SignIn(c echo.Context) error {
	form := new(SignInForm)

	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, err)
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

func CurrentUser(c echo.Context) error {
	user := currentUserByJwtToken(c)
	return c.JSON(http.StatusOK, user)
}

func currentUserByJwtToken(c echo.Context) (user models.User) {
	// Get user id by Jwt token
	result := c.Get("user").(*jwt.Token)
	claims := result.Claims.(jwt.MapClaims)
	id := claims["email"].(string)

	database.Connection().Where("email = ?", id).First(&user)
	return user
}

func generateJwtToken(c echo.Context, user models.User) error {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
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
