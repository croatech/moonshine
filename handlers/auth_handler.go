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

func SignUp(c echo.Context) error {
	db := database.Connect()

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	user := &models.User{
		Username: c.FormValue("username"),
		Email:    c.FormValue("email"),
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

	formEmail := c.FormValue("email")
	formPassword := c.FormValue("password")

	user := models.User{}
	if db.Where("email = ?", formEmail).First(&user).RecordNotFound() {
		return c.JSON(http.StatusInternalServerError, "Email not found or incorrect password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(formPassword))
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
