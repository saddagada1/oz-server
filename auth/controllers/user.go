package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saddagada1/oz-auth/models"
	"github.com/saddagada1/oz-auth/utils"
	"golang.org/x/crypto/bcrypt"
)

// gClient := os.Getenv("GOOGLE_API_CLIENT")
// gSecret := os.Getenv("GOOGLE_API_SECRET")

func Signup(c *gin.Context) {
	var body utils.BasicUserRequest

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "malformed body",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "malformed server: password",
		})
	}

	user := models.User{Email: body.Email, Username: body.Username, Password: string(hash)}
	result := utils.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "malformed server: possible db conflict",
		})

		fmt.Print(result.Error)
	}

	accessToken, refreshToken, err := utils.CreateAuthTokens(user.ID, user.TokenVersion)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "malformed server: tokens",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func Login(c *gin.Context) {
	var body utils.AuthUserRequest

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "malformed body",
		})

		return
	}

	var user models.User
	utils.DB.Where("username = ? OR email = ?", body.Principle).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid login",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid login",
		})
	}

	accessToken, refreshToken, err := utils.CreateAuthTokens(user.ID, user.TokenVersion)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "malformed server: tokens",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"authorization": "ok",
	})
}
