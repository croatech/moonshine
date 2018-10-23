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
	Username string
	Email    string
	Password string
}

type SignInForm struct {
	Email    string
	Password string
}

func SignUp(c echo.Context) error {
	form := SignUpForm{
		Username: c.FormValue("username"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password")}

	user := models.User{
		Username: form.Username,
		Email:    form.Email,
		Password: models.HashPassword(form.Password),
	}

	err := database.Connection().Create(&user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return generateJwtToken(c, user)
	}
}

func SignIn(c echo.Context) error {
	form := SignInForm{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password")}

	user := models.User{}
	if database.Connection().Where("email = ?", form.Email).First(&user).RecordNotFound() {
		return c.JSON(http.StatusInternalServerError, "Email not found or incorrect password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Email not found or incorrect password")
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
