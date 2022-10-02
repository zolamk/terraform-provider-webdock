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
	PhpVersion *string `json:"phpVersion" mapstructure:"php_version"`

	// Image slug
	Slug string `json:"slug,omitempty" mapstructure:"slug"`

	// Web server type
	WebServer *string `json:"webServer" mapstructure:"web_server"`
}

// ServerImages is a collection of ServerImage
type ServerImages []ServerImage

func (c *Client) GetServersImages(ctx context.Context) (*ServerImages, error) {
	req, err := newGetServersImagesRequest(c.Server)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if errorStatus(res.StatusCode) {
		return nil, fmt.Errorf("error getting images: %s", res.Status)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var serverImages ServerImages

	if err = decoder.Decode(&serverImages); err != nil {
		return nil, err
	}

	return &serverImages, nil
}

// newGetServersImagesRequest generates requests for GetServersImages
func newGetServersImagesRequest(server string) (*http.Request, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "images"

	return http.NewRequest("GET", serverURL.String(), nil)
}
