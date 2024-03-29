package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ServerImage model
type ServerImage struct {
	// Image name
	Name string `json:"name,omitempty" mapstructure:"name"`

	// PHP Version. For example &quot;7.4&quot;
	PhpVersion string `json:"phpVersion" mapstructure:"php_version"`

	// Image slug
	Slug string `json:"slug,omitempty" mapstructure:"slug"`

	// Web server type
	WebServer string `json:"webServer" mapstructure:"web_server"`
}

// ServerImages is a collection of ServerImage
type ServerImages []ServerImage

func (c *Client) GetServersImages(ctx context.Context) (ServerImages, error) {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "images"

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error getting server images: %w", err)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting server images: %w", err)
	}

	if errorStatus(res.StatusCode) {
		apiError := APIError{}

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error decoding get server images error response body: %w", err)
		}

		return nil, fmt.Errorf("error getting server images: %w", apiError)
	}

	defer res.Body.Close()

	serverImages := ServerImages{}

	if err = json.NewDecoder(res.Body).Decode(&serverImages); err != nil {
		return nil, fmt.Errorf("error decoding get server images response body: %w", err)
	}

	return serverImages, nil
}
