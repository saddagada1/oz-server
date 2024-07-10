package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(userId uint, version uint) (string, error) {
	expiry, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRES_IN"))

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Duration(expiry)).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
}

func CreateRefreshToken(userId uint, version uint) (string, error) {
	expiry, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRES_IN"))

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Duration(expiry)).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
}

func CreateAuthTokens(userId uint, version uint) (string, string, error) {
	accessToken, err := CreateAccessToken(userId, version)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := CreateRefreshToken(userId, version)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
