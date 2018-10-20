package handlers

import (
	"net/http"
	"sunlight/config/database"
	"sunlight/models"
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
	db := database.Connect()

	form := SignUpForm{
		Username: c.FormValue("username"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password")}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	user := &models.User{
		Username: form.Username,
		Email:    form.Email,
		Password: string(passwordHash),
	}

	err := db.Create(user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, user)
	}
}

func GenerateJwtToken(c echo.Context) error {
	db := database.Connect()

	form := SignInForm{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password")}

	user := models.User{}
	if db.Where("email = ?", form.Email).First(&user).RecordNotFound() {
		return c.JSON(http.StatusInternalServerError, "Email not found or incorrect password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Email not found or incorrect password")
	} else {
		return generateJwtToken(c, user)
	}
}

func generateJwtToken(c echo.Context, user models.User) error {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
