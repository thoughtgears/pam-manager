package handlers

import (
	"net/http"
	"strings"

	"github.com/thoughtgears/pam-manager/models"
	"github.com/thoughtgears/pam-manager/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type PamHandler struct {
	pamService *services.PAMService
}

func NewPamHandler() *PamHandler {
	return &PamHandler{}
}

func (h *PamHandler) GetGrants(c *gin.Context) {
	service, err := services.NewPAMService(c, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create PAM service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PAM service"})
		return
	}

	h.pamService = service

	project := c.Query("project")
	entitlement := c.Query("entitlement")

	if project == "" || entitlement == "" {
		log.Error().Msg("project and entitlement query parameters are required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "project and entitlement query parameters are required"})
		return
	}

	grantsResponse, err := h.pamService.GetGrants(c, project, entitlement)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get grants")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get grants"})
		return
	}

	var grants []models.Grant

	for _, grant := range grantsResponse {
		var id string
		parts := strings.Split(grant.Name, "/")
		if len(parts) > 0 {
			id = parts[len(parts)-1]
		}

		g := models.Grant{
			ID:            id,
			Name:          grant.Name,
			Requester:     grant.Requester,
			Duration:      grant.RequestedDuration.GetSeconds(),
			Justification: grant.Justification.GetUnstructuredJustification(),
			State:         grant.State.String(),
		}

		for _, role := range grant.PrivilegedAccess.GetGcpIamAccess().RoleBindings {
			g.Roles = append(g.Roles, role.Role)
		}

		grants = append(grants, g)
	}

	c.JSON(http.StatusOK, gin.H{"grants": grants})
}

func (h *PamHandler) RequestGrant(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	authToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	token := &oauth2.Token{
		AccessToken: authToken,
		TokenType:   "Bearer",
	}

	tokenSource := oauth2.StaticTokenSource(token)
	service, err := services.NewPAMService(c, tokenSource)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create PAM service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PAM service"})
		return
	}

	h.pamService = service

	var req struct {
		ProjectID   string `json:"project_id" binding:"required"`
		Entitlement string `json:"entitlement" binding:"required"`
		Reason      string `json:"reason" binding:"required"`
		Duration    int64  `json:"duration" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	grantResponse, err := h.pamService.RequestGrant(c, req.ProjectID, req.Entitlement, req.Reason, req.Duration)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create grant")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create grant"})
		return
	}
	var id string
	parts := strings.Split(grantResponse.Name, "/")
	if len(parts) > 0 {
		id = parts[len(parts)-1]
	}

	grant := models.Grant{
		ID:            id,
		Name:          grantResponse.Name,
		Requester:     grantResponse.Requester,
		Duration:      grantResponse.RequestedDuration.GetSeconds(),
		Justification: grantResponse.Justification.GetUnstructuredJustification(),
		State:         grantResponse.State.String(),
	}

	for _, role := range grantResponse.PrivilegedAccess.GetGcpIamAccess().RoleBindings {
		grant.Roles = append(grant.Roles, role.Role)
	}

	c.JSON(http.StatusOK, gin.H{"grant": grant})
}

func (h *PamHandler) ApproveGrant(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	authToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	token := &oauth2.Token{
		AccessToken: authToken,
		TokenType:   "Bearer",
	}

	tokenSource := oauth2.StaticTokenSource(token)
	service, err := services.NewPAMService(c, tokenSource)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create PAM service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PAM service"})
		return
	}

	h.pamService = service

	id := c.Param("id")

	var req struct {
		ProjectID   string `json:"project_id" binding:"required"`
		Entitlement string `json:"entitlement" binding:"required"`
		Reason      string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	grantResponse, err := h.pamService.ApproveGrant(c, id, req.ProjectID, req.Entitlement, req.Reason)
	if err != nil {
		log.Error().Err(err).Msg("Failed to approve grant")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve grant"})
		return
	}

	grant := models.Grant{
		ID:            id,
		Name:          grantResponse.Name,
		Requester:     grantResponse.Requester,
		Duration:      grantResponse.RequestedDuration.GetSeconds(),
		Justification: grantResponse.Justification.GetUnstructuredJustification(),
		State:         grantResponse.State.String(),
	}

	for _, role := range grantResponse.PrivilegedAccess.GetGcpIamAccess().RoleBindings {
		grant.Roles = append(grant.Roles, role.Role)
	}

	c.JSON(http.StatusOK, gin.H{"grant": grant})
}

func (h *PamHandler) RevokeGrant(c *gin.Context) {
	service, err := services.NewPAMService(c, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create PAM service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PAM service"})
		return
	}

	h.pamService = service

	id := c.Param("id")
	project := c.Query("project")
	entitlement := c.Query("entitlement")
	reason := c.Query("reason")

	if id == "" {
		log.Error().Msg("id parameter is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	if project == "" || entitlement == "" {
		log.Error().Msg("project and entitlement query parameters are required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "project and entitlement query parameters are required"})
		return
	}

	if reason == "" {
		reason = "Automated revocation, no reason provided"
	}

	grantResponse, err := h.pamService.RevokeGrant(c, id, project, entitlement, reason)
	if err != nil {
		log.Error().Err(err).Msg("Failed to revoke grant")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke grant"})
		return
	}

	grant := models.Grant{
		ID:            id,
		Name:          grantResponse.Name,
		Requester:     grantResponse.Requester,
		Duration:      grantResponse.RequestedDuration.GetSeconds(),
		Justification: grantResponse.Justification.GetUnstructuredJustification(),
		State:         grantResponse.State.String(),
	}

	for _, role := range grantResponse.PrivilegedAccess.GetGcpIamAccess().RoleBindings {
		grant.Roles = append(grant.Roles, role.Role)
	}

	c.JSON(http.StatusOK, gin.H{"grant": grantResponse})
}
