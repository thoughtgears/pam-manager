package handlers

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"golang.org/x/oauth2"

	"github.com/thoughtgears/pam-manager/services"

	"github.com/gin-gonic/gin"
)

type PAMHandler struct {
	pamService *services.PAMService
}

func NewPAMHandler(pamService *services.PAMService) *PAMHandler {
	return &PAMHandler{pamService: pamService}
}

func (h *PAMHandler) RequestGrant(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	authToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	token := &oauth2.Token{
		AccessToken: authToken,
		TokenType:   "Bearer",
	}

	tokenSource := oauth2.StaticTokenSource(token)

	var req struct {
		ProjectID     string `json:"projectId" binding:"required"`
		Entitlement   string `json:"entitlement" binding:"required"`
		Justification string `json:"justification" binding:"required"`
		Duration      int64  `json:"duration" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	grant, err := h.pamService.RequestGrant(ctx, tokenSource, req.ProjectID, req.Entitlement, req.Justification, req.Duration)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create grant")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create grant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"grant": grant})
}
