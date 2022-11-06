package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ServerLocation model
type ServerLocation struct {
	// Location City
	City string `json:"city,omitempty" mapstructure:"city"`

	// Location Country
	Country string `json:"country,omitempty" mapstructure:"country"`

	// Location Description
	Description string `json:"description,omitempty" mapstructure:"description"`

	// Location Icon
	Icon string `json:"icon,omitempty" mapstructure:"icon"`

	// Location ID
	ID string `json:"id,omitempty" mapstructure:"id"`

	// Location Name
	Name string `json:"name,omitempty" mapstructure:"name"`
}

type ServerLocations []ServerLocation

func (c *Client) GetServersLocations(ctx context.Context) (ServerLocations, error) {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "locations"

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		apiError := APIError{}

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error decoding get server locations error response body: %w", err)
		}

		return nil, fmt.Errorf("error getting server locations: %w", apiError)
	}

	locations := ServerLocations{}

	if err = json.NewDecoder(res.Body).Decode(&locations); err != nil {
		return nil, fmt.Errorf("error decoding get server locations response body: %w", err)
	}

	return locations, nil
}
