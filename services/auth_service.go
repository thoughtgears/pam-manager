package services

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	oauthConfig *oauth2.Config
}

func NewAuthService(clientID, clientSecret, redirectURL string) *AuthService {
	return &AuthService{
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/cloud-platform"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (a *AuthService) GetLoginURL(state string) string {
	return a.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// HandleCallback exchanges the authorization code for a token and returns both the token and the HTTP client
func (a *AuthService) HandleCallback(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := a.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}
