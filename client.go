package netatmo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/exzz/netatmo-api-go/types"
)

type OAuthScope string

const (
	OAuthScopeWeather OAuthScope = "read_station"
	OAuthScopeEnergy  OAuthScope = "read_thermostat"
)

var (
	// ErrNotAuthenticated is returned from the client when it is not authenticated yet.
	ErrNotAuthenticated = errors.New("no token available")

	defaultScopes = []string{
		string(OAuthScopeWeather),
	}
)

// TokenUpdateFunc defines a function that can act as a callback for a token update.
type TokenUpdateFunc func(new *oauth2.Token)

// Config is used to specify credential to Netatmo API
type Config struct {
	// ClientID from netatmo app registration at http://dev.netatmo.com/dev/listapps
	ClientID string
	// ClientSecret Client app secret
	ClientSecret string
	// Scopes contains the scopes that should be used for authentication
	Scopes []OAuthScope
}

// Client use to make request to Netatmo API
type Client struct {
	oauth          *oauth2.Config
	httpClient     *http.Client
	updateCallback TokenUpdateFunc
}

// NewClient creates an unauthenticated NetAtmo API client.
func NewClient(config Config, tokenCallback TokenUpdateFunc) *Client {
	scopes := []string{}
	for _, scope := range config.Scopes {
		scopes = append(scopes, string(scope))
	}

	if len(scopes) == 0 {
		scopes = defaultScopes
	}

	oauth := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	return &Client{
		oauth:          oauth,
		updateCallback: tokenCallback,
	}
}

// AuthCodeURL creates an authentication URL that can be passed to the user.
func (c *Client) AuthCodeURL(redirectURL, state string) string {
	c.oauth.RedirectURL = redirectURL
	return c.oauth.AuthCodeURL(state)
}

// Exchange converts an authentication code into a token and authenticates the client.
func (c *Client) Exchange(ctx context.Context, code, state string) error {
	token, err := c.oauth.Exchange(ctx, code, oauth2.SetAuthURLParam("state", state))
	if err != nil {
		return err
	}

	c.InitWithToken(ctx, token)
	return nil
}

// CurrentToken retrieves the token for persisting state.
func (c *Client) CurrentToken() (*oauth2.Token, error) {
	if c.httpClient == nil {
		return nil, ErrNotAuthenticated
	}

	transport := c.httpClient.Transport.(*oauth2.Transport)
	source := transport.Source
	return source.Token()
}

func (c *Client) tokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	source := c.oauth.TokenSource(ctx, token)
	if c.updateCallback == nil {
		return source
	}

	return &callbackTokenSource{
		callback:    c.updateCallback,
		tokenSource: c.oauth.TokenSource(ctx, token),
		lastToken:   token,
	}
}

// InitWithToken initializes the client with an existing token.
func (c *Client) InitWithToken(ctx context.Context, token *oauth2.Token) {
	c.httpClient = oauth2.NewClient(ctx, c.tokenSource(ctx, token))
}

func handleError(resp *http.Response) error {
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return fmt.Errorf("error reading body for status code %d: %w", resp.StatusCode, err)
	}

	var errResp types.ErrorResponse
	if err := json.Unmarshal(buf.Bytes(), &errResp); err != nil {
		return fmt.Errorf("can not parse error message for status %d: %s - parse error: %w", resp.StatusCode, buf.String(), err)
	}

	if errResp.Error.Message != "" {
		return fmt.Errorf("got error %d: %s (HTTP status %d)", errResp.Error.Code, errResp.Error.Message, resp.StatusCode)
	}

	return fmt.Errorf("got non-ok HTTP status %d: %s", resp.StatusCode, buf.String())
}
