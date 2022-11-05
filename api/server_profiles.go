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

func (c *Client) GetServersProfiles(ctx context.Context, params *GetServersProfilesParams) (ServerProfiles, error) {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "profiles"

	queryValues, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	serverURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	profiles := ServerProfiles{}

	if err := json.NewDecoder(res.Body).Decode(&profiles); err != nil {
		return nil, err
	}

	return profiles, err
}
