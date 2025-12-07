package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	ErrFailedToExchange    = errors.New("failed to exchange authorization code")
	ErrFailedToGetUserInfo = errors.New("failed to get user info from Google")
)

// GoogleUserInfo represents the user info from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// Config holds OAuth configuration
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// GoogleOAuth handles Google OAuth2 operations
type GoogleOAuth struct {
	config *oauth2.Config
}

// NewGoogleOAuth creates a new Google OAuth handler
func NewGoogleOAuth(cfg Config) *GoogleOAuth {
	return &GoogleOAuth{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// GetAuthURL returns the URL to redirect the user for authentication
func (g *GoogleOAuth) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Exchange exchanges the authorization code for tokens
func (g *GoogleOAuth) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, ErrFailedToExchange
	}
	return token, nil
}

// GetUserInfo retrieves user info from Google using the access token
func (g *GoogleOAuth) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := g.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, ErrFailedToGetUserInfo
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrFailedToGetUserInfo
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrFailedToGetUserInfo
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, ErrFailedToGetUserInfo
	}

	return &userInfo, nil
}
