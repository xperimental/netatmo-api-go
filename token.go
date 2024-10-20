package netatmo

import (
	"sync"

	"golang.org/x/oauth2"
)

type callbackTokenSource struct {
	callback    TokenUpdateFunc
	tokenSource oauth2.TokenSource
	lastToken   *oauth2.Token
	sync.Mutex
}

func (c *callbackTokenSource) Token() (*oauth2.Token, error) {
	c.Lock()
	defer c.Unlock()

	token, err := c.tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if tokensDiffer(c.lastToken, token) {
		c.callback(token)
		c.lastToken = token
	}

	return token, nil
}

func tokensDiffer(old, new *oauth2.Token) bool {
	if old == nil && new == nil {
		return false
	}

	if old == nil || new == nil {
		return true
	}

	return old.RefreshToken != new.RefreshToken ||
		old.AccessToken != new.AccessToken ||
		old.Expiry != new.Expiry
}
