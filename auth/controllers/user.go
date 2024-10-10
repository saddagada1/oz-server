package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/saddagada1/oz-auth/models"
	"github.com/saddagada1/oz-auth/utils"
	"golang.org/x/crypto/bcrypt"
)

// gClient := os.Getenv("GOOGLE_API_CLIENT")
// gSecret := os.Getenv("GOOGLE_API_SECRET")

func parseDuplicatedKeyError(err error) (string, string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		constraint := pgErr.ConstraintName
		if strings.Contains(constraint, "email") {
			return "email", constraint
		} else if strings.Contains(constraint, "username") {
			return "username", constraint
		}
	}
	return "", ""
}

func Signup(c *gin.Context) {
	var body utils.BasicUserRequest

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"subject": "request",
			"message": "malformed body",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"subject": "malformed server",
			"message": "failed to hash password",
		})

		return
	}

	user := models.User{Email: body.Email, Username: body.Username, Password: string(hash)}
	result := utils.DB.Create(&user)

	if result.Error != nil {
		field, duplicatedKey := parseDuplicatedKeyError(result.Error)

		if duplicatedKey != "" {
			if field == "email" {
				c.JSON(http.StatusConflict, gin.H{
					"subject": "email",
					"message": "conflict",
				})
			} else if field == "username" {
				c.JSON(http.StatusConflict, gin.H{
					"subject": "username",
					"message": "conflict",
				})
			}
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"subject": "malformed server",
			"message": "possible db conflict or fatal :(",
		})

		return
	}

	accessToken, refreshToken, err := utils.CreateAuthTokens(user.ID, user.TokenVersion)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"subject": "malformed server",
			"message": "failed to generate tokens",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         user,
	})
}

func Login(c *gin.Context) {
	var body utils.AuthUserRequest

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"subject": "request",
			"message": "malformed body",
		})

		return
	}

	var user models.User
	utils.DB.Where("username = ? OR email = ?", body.Principle, body.Principle).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"subject": "principle",
			"message": "invalid login",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"subject": "principle",
			"message": "invalid login",
		})

		return
	}

	accessToken, refreshToken, err := utils.CreateAuthTokens(user.ID, user.TokenVersion)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"subject": "malformed server",
			"message": "failed to generate tokens",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         user,
	})
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func RefreshToken(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(http.StatusBadRequest, gin.H{"subject": "request", "message": "no auth header"})
		return
	}

	if !strings.HasPrefix(header, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{"subject": "request", "message": "malformed auth header"})
		return
	}

	token := strings.TrimPrefix(header, "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"subject": "token", "message": "malformed token"})
		return
	}

	payload, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signature: %v", t.Header)
		}

		return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"subject": "token", "message": "invalid token"})
		return
	}

	if claims, ok := payload.Claims.(jwt.MapClaims); ok && payload.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"subject": "token", "message": "expired token"})
			return
		}

		var user models.User
		utils.DB.Where("ID = ?", claims["sub"]).First(&user)

		if user.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"subject": "user",
				"message": "user not found",
			})

			return
		}

		accessToken, refreshToken, err := utils.CreateAuthTokens(user.ID, user.TokenVersion)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"subject": "malformed server",
				"message": "failed to generate tokens",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"user":         user,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"subject": "token", "message": "invalid token"})
	}
}
