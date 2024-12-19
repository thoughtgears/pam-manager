package handlers

import (
	"encoding/json"
	"fmt"
	"io"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

func Debug(c *gin.Context) {
	log.Info().Msg("--- Request Headers ---")
	for key, values := range c.Request.Header {
		log.Info().Str("header", fmt.Sprintf("%s: %s\n", key, strings.Join(values, ", "))).Msg("")
	}

	log.Info().Msg("--- Request Body ---")
	body, err := c.GetRawData()
	if err != nil {
		log.Error().Err(err).Msg("Error reading request body")
	} else {
		log.Info().Str("body", string(body)).Msg("")
	}

	// Restore the body for future handlers
	c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Error().Msg("No Authorization header found")
	} else {
		log.Info().Str("Authorization", authHeader).Msg("Authorization header found")

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			log.Error().Err(err).Msg("Error parsing JWT token")
		} else {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				claimsJSON, _ := json.MarshalIndent(claims, "", "  ")
				log.Info().Str("claims", string(claimsJSON)).Msg("JWT claims")
			} else {
				log.Error().Msg("Error parsing JWT claims")
			}
		}
	}

	c.JSON(200, gin.H{"message": "Debug complete"})
}
