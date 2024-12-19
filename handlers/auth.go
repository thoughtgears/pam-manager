package handlers

import (
	"net/http"

	"github.com/thoughtgears/pam-manager/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	url := h.authService.GetLoginURL("state-token")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	token, err := h.authService.HandleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Return relevant token information
	c.JSON(http.StatusOK, gin.H{
		"message":     "Authentication successful",
		"accessToken": token.AccessToken,
		"expiry":      token.Expiry,
	})
}
