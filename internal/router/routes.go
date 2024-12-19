package router

import (
	"github.com/thoughtgears/pam-manager/handlers"
	"github.com/thoughtgears/pam-manager/internal/router/middleware"

	"github.com/gin-gonic/gin"
)

func (r *Router) RegisterRoutes(authHandler *handlers.AuthHandler, pamHandler *handlers.PamHandler) {
	r.engine.Use(gin.Recovery(), middleware.Logger())

	r.engine.POST("/debug", middleware.AuthRequired(), handlers.Debug)

	// Auth routes
	auth := r.engine.Group("/auth")
	{
		auth.GET("/google/login", authHandler.Login)
		auth.GET("/google/callback", authHandler.Callback)
	}

	// PAM routes
	pam := r.engine.Group("/pam")
	pam.Use(middleware.AuthRequired())
	{
		pam.GET("/grants", pamHandler.GetGrants)
		pam.POST("/grants", pamHandler.RequestGrant)
		pam.PATCH("/grants/:id", pamHandler.ApproveGrant)
		pam.DELETE("/grants/:id", pamHandler.RevokeGrant)
	}
}
