package types

import (
	"time"

	"golang.org/x/oauth2"
)

// ErrorResponse contains information about an error.
type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// TokenResponse contains the authentication token received from the API
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (t TokenResponse) Token(issueTime time.Time) *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  t.AccessToken,
		TokenType:    "",
		RefreshToken: t.RefreshToken,
		Expiry:       issueTime.Add(time.Second * time.Duration(t.ExpiresIn)),
	}
}
