package api

import (
	"context"
	"encoding/json"
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
	Id string `json:"id,omitempty" mapstructure:"id"`

	// Location Name
	Name string `json:"name,omitempty" mapstructure:"name"`
}

type ServerLocations []ServerLocation

func (c *Client) GetServersLocations(ctx context.Context, reqEditors ...RequestEditorFn) (*ServerLocations, error) {
	req, err := newGetServersLocationsRequest(c.Server)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var locations ServerLocations

	if err = decoder.Decode(&locations); err != nil {
		return nil, err
	}

	return &locations, nil
}

// newGetServersLocationsRequest generates requests for GetServersLocations
func newGetServersLocationsRequest(server string) (*http.Request, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "locations"

	return http.NewRequest("GET", serverURL.String(), nil)
}
