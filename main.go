package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/thoughtgears/pam-manager/handlers"
	"github.com/thoughtgears/pam-manager/internal/config"
	"github.com/thoughtgears/pam-manager/internal/router"
	"github.com/thoughtgears/pam-manager/services"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfg config.Config

func init() {
	zerolog.LevelFieldName = "severity"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	envconfig.MustProcess("", &cfg)
}

func main() {
	authService := services.NewAuthService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)
	authHandler := handlers.NewAuthHandler(authService)
	pamHandler := handlers.NewPamHandler()

	// Create the router
	r, err := router.New(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create router")
	}

	r.RegisterRoutes(authHandler, pamHandler)
	log.Fatal().Err(r.Run()).Msg("Failed to start server")
}
