package config

// Config is the configuration for the application
// It contains the port, debug, and host configuration
type Config struct {
	Port               string `envconfig:"PORT" default:"8080"`
	Debug              bool   `envconfig:"DEBUG" default:"false"`
	SlackClientSecret  string `envconfig:"SLACK_CLIENT_SECRET" required:"true"`
	SlackSigningSecret string `envconfig:"SLACK_SIGNING_SECRET" required:"true"`
	SlackToken         string `envconfig:"SLACK_BOT_TOKEN" required:"true"`
	GoogleClientID     string `envconfig:"GOOGLE_CLIENT_ID" required:"true"`
	GoogleClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET" required:"true"`
	GoogleRedirectURL  string `envconfig:"GOOGLE_REDIRECT_URL" required:"true"`
}
