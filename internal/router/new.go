package router

import (
	"fmt"

	"github.com/thoughtgears/pam-manager/internal/config"
	"github.com/thoughtgears/pam-manager/internal/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Router struct {
	engine *gin.Engine
	debug  bool
	host   string
	port   string
}

// New creates a new Router with the given debug flag
// If debug is true, the router will listen on
// 127.0.0.1 and set the trusted proxies to nil
// If debug is false, the router will listen on
// 0.0.0.0 and set gin to release mode
func New(config *config.Config) (*Router, error) {
	var router Router

	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router.host = "0.0.0.0"
	router.port = config.Port

	router.engine = gin.New()
	router.engine.Use(gin.Recovery(), middleware.Logger())

	if config.Debug {
		if err := router.engine.SetTrustedProxies(nil); err != nil {
			return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
		}
		router.host = "127.0.0.1"
	}

	return &router, nil
}

// Run starts the server on the configured host and port
func (r *Router) Run() error {
	log.Info().Msgf("Starting server on %s:%s", r.host, r.port)
	return r.engine.Run(fmt.Sprintf("%s:%s", r.host, r.port))
}
