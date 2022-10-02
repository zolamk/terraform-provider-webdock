package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

// CPU model
type CPUDTO struct {
	// Number of cores
	Cores int64 `json:"cores,omitempty" mapstructure:"cores"`

	// Number of threads
	Threads int64 `json:"threads,omitempty" mapstructure:"threads"`
}

// GetServersProfilesParams defines parameters for GetServersProfiles.
type GetServersProfilesParams struct {
	// Location of the profile
	LocationId string `json:"locationId,omitempty" url:"locationId,omitempty"`
}

// ServerProfile model
type ServerProfile struct {
	// CPU model
	CPU CPUDTO `json:"cpu,omitempty" mapstructure:"cpu"`

	// Disk size (in MiB)
	Disk int64 `json:"disk,omitempty" mapstructure:"disk"`

	// Profile name
	Name string `json:"name,omitempty" mapstructure:"name"`

	// Price model
	Price PriceDTO `json:"price,omitempty" mapstructure:"-"`

	// RAM memory (in MiB)
	RAM int64 `json:"ram,omitempty" mapstructure:"ram"`

	// Profile slug
	Slug string `json:"slug,omitempty" mapstructure:"slug"`
}

// ServerProfiles is a collection of ServerProfile
type ServerProfiles []ServerProfile

func (c *Client) GetServersProfiles(ctx context.Context, params *GetServersProfilesParams) (*ServerProfiles, error) {
	req, err := newGetServersProfilesRequest(c.Server, params)
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

	var profiles ServerProfiles

	if err := decoder.Decode(&profiles); err != nil {
		return nil, err
	}

	return &profiles, err
}

// newGetServersProfilesRequest generates requests for GetServersProfiles
func newGetServersProfilesRequest(server string, params *GetServersProfilesParams) (*http.Request, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "profiles"

	queryValues, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	serverURL.RawQuery = queryValues.Encode()

	return http.NewRequest("GET", serverURL.String(), nil)
}
