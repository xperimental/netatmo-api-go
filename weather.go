package netatmo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Read returns the list of stations owned by the user and their modules
func (c *Client) Read() (*DeviceCollection, error) {
	if c.httpClient == nil {
		return nil, ErrNotAuthenticated
	}

	data := url.Values{"app_type": {"app_station"}}

	req, err := http.NewRequest("GET", deviceURL, nil)
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
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, resp.Body); err != nil {
			return nil, fmt.Errorf("error reading body for status code %d: %w", resp.StatusCode, err)
		}

		var errResp ErrorResponse
		if err := json.Unmarshal(buf.Bytes(), &errResp); err != nil {
			return nil, fmt.Errorf("can not parse error message for status %d: %s - parse error: %w", resp.StatusCode, buf.String(), err)
		}

		if errResp.Error.Message != "" {
			return nil, fmt.Errorf("got error %d: %s (HTTP status %d)", errResp.Error.Code, errResp.Error.Message, resp.StatusCode)
		}

		return nil, fmt.Errorf("got non-ok HTTP status %d: %s", resp.StatusCode, buf.String())
	}

	result := &DeviceCollection{}
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
