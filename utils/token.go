package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Access token expiration time
	DefaultAccessTokenExpiry = time.Minute * 15

	// Refresh token expiration time
	DefaultRefreshTokenExpiry = time.Hour * 24
)

func GetAccessTokenExpiry() time.Duration {
	// expiry := config.Config.GetInt("ACCESS_TOKEN_EXPIRY_DURATION")
	// if expiry == 0 {
	return DefaultAccessTokenExpiry
	// }
	// return time.Duration(expiry) * time.Second
}

func GetRefreshTokenExpiry() time.Duration {
	// expiry := config.Config.GetInt("REFRESH_TOKEN_EXPIRY_DURATION")
	// if expiry == 0 {
	return DefaultRefreshTokenExpiry
	// }
	// return time.Duration(expiry) * time.Second
}

func GenerateTokens(userID string) (string, string, error) {
	// Secret keys for access and refresh tokens
	accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")

	now := time.Now()

	// Generate access token with short expiration
	accessTokenString, err := GenerateJWTToken(userID, accessSecret, GetAccessTokenExpiry(), now)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token with longer expiration
	refreshTokenString, err := GenerateJWTToken(userID, refreshSecret, GetRefreshTokenExpiry(), now)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func GetForgotPasswordTokenExpiry() time.Duration {
	// expiry := config.Config.GetInt("PASSWORD_TOKEN_EXPIRY_DURATION")
	// if expiry == 0 {
	return 5 * time.Minute
	// }
	// return time.Duration(expiry) * time.Second
}

func GenerateJWTToken(sub interface{}, secret string, expiry time.Duration, iat ...time.Time) (string, error) {
	now := time.Now()
	if len(iat) > 0 {
		now = iat[0]
	}
	// Generate access token with short expiration
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"iat": now.Unix(),
		"exp": now.Add(expiry).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(token, secret string) (jwt.MapClaims, error) {
	// Validate and parse refresh token
	refreshToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Ensure signature method and validity
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		if _, ok := token.Claims.(jwt.MapClaims); !ok {
			return nil, fmt.Errorf("invalid claims")
		}
		return []byte(secret), nil
	})
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("invalid token")
	}

	// fetch user ID from refresh token claims
	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || !refreshToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
