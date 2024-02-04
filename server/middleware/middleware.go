package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type FailedResponse struct {
	Error string `json:"error"`
}

// AuthenticationMiddleware checks if the user has a valid JWT token
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, FailedResponse{Error: "authentication failed"})
			c.Abort()
			return
		}

		// The token should be prefixed with "Bearer "
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, FailedResponse{Error: "authentication failed"})
			c.Abort()
			return
		}

		tokenString = tokenParts[1]

		claims, err := VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, FailedResponse{Error: "authentication failed"})
			c.Abort()
			return
		}

		c.Set("sub", claims["sub"])
		c.Next()
	}
}

// VerifyToken verifies a token JWT validate
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC with SHA-256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})

	// Check for errors
	if err != nil {
		return nil, err
	}

	// Validate the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}
