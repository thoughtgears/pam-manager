package main

import (
	"net/http"

	"github.com/thoughtgears/pam-manager/internal/middleware"

	"cloud.google.com/go/privilegedaccessmanager/apiv1/privilegedaccessmanagerpb"
	cloudevent "github.com/cloudevents/sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

// Config is the configuration for the application
// It contains the port, debug, and host configuration
type Config struct {
	Port  string `envconfig:"PORT" default:"8080"`
	Debug bool   `envconfig:"DEBUG" default:"false"`
	Host  string `envconfig:"HOST" default:"0.0.0.0"`
}

var config Config

// init sets up the logging configuration
// and processes the environment variables
func init() {
	zerolog.LevelFieldName = "severity"
	envconfig.MustProcess("", &config)
}

// main is the entrypoint for the application
// It sets up the Gin router and starts the server
// It also sets the log level based on the DEBUG environment variable
// and sets the trusted proxies based on the DEBUG environment variable
// It also sets up the health check endpoint
// and the endpoint for handling grant requests
func main() {
	router := gin.New()
	router.Use(gin.Recovery(), middleware.Logger())

	switch config.Debug {
	case true:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		config.Host = "127.0.0.1"
		if err := router.SetTrustedProxies(nil); err != nil {
			log.Fatal().Err(err).Msg("Failed to set trusted proxies")
		}
	case false:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		log.Fatal().Msg("Debug environment variable must be set to true or false")
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/", GrantRequest())

	log.Info().Msgf("Starting server on %s:%s", config.Host, config.Port)
	log.Fatal().Err(router.Run(config.Host + ":" + config.Port)).Msg("Failed to start server")
}

// GrantRequest is a Gin handler that processes a CloudEvent
// containing a privilegedaccessmanagerpb.CreateGrantRequest
// and logs the request.
func GrantRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		event, err := cloudevent.NewEventFromHTTPRequest(ctx.Request)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create CloudEvent from request")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create CloudEvent from request"})
			return
		}

		var data privilegedaccessmanagerpb.CreateGrantRequest
		umo := &protojson.UnmarshalOptions{DiscardUnknown: true}

		// Unmarshal CloudEvent data into the request struct
		if err := umo.Unmarshal(event.Data(), &data); err != nil {
			log.Error().Err(err).Msg("Failed to parse grant request")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse grant request"})
			return
		}

		// Safely marshal data into JSON for logging
		jsonData, err := protojson.Marshal(&data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal grant request")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}
		log.Info().Msgf("Received grant request: %s", string(jsonData))
		ctx.Status(http.StatusAccepted)
	}
}
