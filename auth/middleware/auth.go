package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no auth header"})
		c.Abort()
		return
	}

	if !strings.HasPrefix(header, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed auth header"})
		c.Abort()
		return
	}

	token := strings.TrimPrefix(header, "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "malformed token"})
		c.Abort()
		return
	}

	payload, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signature: %v", t.Header)
		}

		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})

	if claims, ok := payload.Claims.(jwt.MapClaims); ok && payload.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "expired token"})
			c.Abort()
			return
		}

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
