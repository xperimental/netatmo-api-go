package netatmo

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/exzz/netatmo-api-go/types/energy"
)

func (c *Client) GetHomes() ([]energy.Home, error) {
	if c.httpClient == nil {
		return nil, ErrNotAuthenticated
	}

	req, err := http.NewRequest("GET", homesURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleError(resp)
	}

	result := &energy.HomesDataResponse{}
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}

	return result.Body.Homes, nil
}

func (c *Client) GetHomeStatus(homeID string) (*energy.Home, error) {
	if c.httpClient == nil {
		return nil, ErrNotAuthenticated
	}

	data := url.Values{
		"home_id": []string{homeID},
	}

	req, err := http.NewRequest("GET", homeStatusURL, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = data.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleError(resp)
	}

	result := &energy.HomeStatusResponse{}
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}

	return &result.Body.Home, nil
}
