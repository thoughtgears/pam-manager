package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
)

const UserContextKey = "user"

// AuthRequired middleware checks for a valid OAuth2 token
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		oauth2Service, err := oauth2.NewService(context.Background(), option.WithoutAuthentication())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize OAuth2 service"})
			return
		}

		tokenInfo, err := oauth2Service.Tokeninfo().AccessToken(tokenString).Do()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set(UserContextKey, tokenInfo)
		c.Next()
	}
}
