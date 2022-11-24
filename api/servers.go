package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

// Charge summary items model
type ChargeSummaryItem struct {
	// Price model
	Price Price `json:"price,omitempty"`

	// Charge text
	Text string `json:"text,omitempty"`
}

// Server resize model
type ChargeSummary struct {
	// True if the amount will be refunded
	IsRefund bool `json:"isRefund,omitempty"`

	// List of charges to be applied or refunded
	Items []ChargeSummaryItem `json:"items,omitempty"`

	// Server resize model
	Total ChargeSummaryTotal `json:"total,omitempty"`
}

// PatchServer model
type PatchServerRequestBody struct {
	// Name of the server
	Name string `json:"name"`
}

// Reinstall Server model
type ReinstallServerRequestBody struct {
	// Image slug of the image you want to reload the server with. Any image listed for the server location from /images is valid.
	ImageSlug string `json:"imageSlug"`
}

// Resize server model
type ResizeServerRequestBody struct {
	// Profile slug to resize to
	ProfileSlug string `json:"profileSlug"`
}

// GetServersParams defines parameters for GetServers.
type GetServersParams struct {
	// Filter by current status of the server
	Status string `form:"status,omitempty" json:"status,omitempty"`
}

// Warning model
type Warning struct {
	// Warning message
	Data map[string]interface{} `json:"data,omitempty"`

	// Warning message
	Message string `json:"message,omitempty"`

	// Warning type
	Type string `json:"type,omitempty"`
}

// Server model
type Server struct {
	// SSH Password Authentication Enabled for this Server
	SSHPasswordAuthEnabled bool `json:"SSHPasswordAuthEnabled,omitempty" mapstructure:"ssh_password_auth_enabled"`

	// Wordpress lockdown status
	WordPressLockDown bool `json:"WordPressLockDown,omitempty" mapstructure:"wordpress_lockdown"`

	// Aliases - Domain names for the server as known by Webdock. First entry should be treated as the &quot;Main Domain&quot; for the server.
	Aliases []string `json:"aliases,omitempty" mapstructure:"aliases"`

	// Creation date/time
	Date string `json:"date,omitempty" mapstructure:"created_at"`

	// Server image
	Image string `json:"image,omitempty" mapstructure:"image_slug"`

	// IPv4 address
	Ipv4 string `json:"ipv4" mapstructure:"ipv4"`

	// IPv6 address
	Ipv6 string `json:"ipv6" mapstructure:"ipv6"`

	// Location ID of the server
	Location string `json:"location" mapstructure:"location_id"`

	// Server name
	Name string `json:"name,omitempty" mapstructure:"name"`

	// Server profile
	Profile string `json:"profile" mapstructure:"profile_slug"`

	// Server slug
	Slug string `json:"slug,omitempty" mapstructure:"slug"`

	// Last known snapshot runtime (seconds)
	SnapshotRunTime int64 `json:"snapshotRunTime,omitempty" mapstructure:"snapshot_runtime"`

	// Server status
	Status string `json:"status,omitempty" mapstructure:"status"`

	// Server virtualization type indicating whether it's a Webdock LXD VPS or a KVM Virtual Machine
	Virtualization string `json:"virtualization,omitempty" mapstructure:"virtualization"`

	// Webserver type
	WebServer string `json:"webServer,omitempty" mapstructure:"webserver"`

	CallbackID string `json:"-" mapstructure:"-"`
}

// Servers is a collection of Server
type Servers []Server

// Server resize model
type ServerResize struct {
	ChargeSummary *ChargeSummary `json:"chargeSummary,omitempty"`
	Warnings      []Warning      `json:"warnings,omitempty"`
}

// Server resize model
type ChargeSummaryTotal struct {
	// Price model
	SubTotal Price `json:"subTotal,omitempty"`

	// Price model
	Total Price `json:"total,omitempty"`

	// Price model
	Vat Price `json:"vat,omitempty"`
}

// Post Server model
type CreateServerRequestBody struct {
	// Slug of the server image. Get this from the /images endpoint. You must pass either this parameter or snapshotId
	ImageSlug string `json:"imageSlug,omitempty"`

	// ID of the location. Get this from the /locations endpoint.
	LocationId string `json:"locationId"`

	// Name of the server
	Name string `json:"name"`

	// Slug of the server profile. Get this from the /profiles endpoint.
	ProfileSlug string `json:"profileSlug"`

	// Suggested Slug (shortname) of the server. Up to 12 alphanumeric chars. This slug is effectively your server ID and anything submitted in this field is merely your suggestion for a slug. If omitted or the suggested slug is already taken, the system will automatically generate an appropriate unique slug based on your server Name or suggestion. Always check the return from this method to determine the actual slug your server ended up receiving.
	Slug string `json:"slug,omitempty"`

	// SnapshotID from which to create the server. Get this from the /servers/{serverSlug}/snapshots endpoint. You must pass either this parameter or imageSlug.
	SnapshotId int64 `json:"snapshotId,omitempty"`

	// Virtualization type for your new server. container means the server will be a Webdock LXD VPS and kvm means it will be a KVM Virtual machine. If you specify a snapshotId in the request, the server type from which the snapshot belongs much match the virtualization selected. Reason being that KVM images are incompatible with LXD images and vice-versa.
	Virtualization string `json:"virtualization,omitempty"`
}

func (c *Client) GetServers(ctx context.Context, params GetServersParams) (Servers, error) {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "servers"

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

	if errorStatus(res.StatusCode) {
		apiError := APIError{}

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error decoding get servers error response body: %w", err)
		}

		return nil, fmt.Errorf("error getting servers: %w", apiError)
	}

	servers := Servers{}

	if err = json.NewDecoder(res.Body).Decode(&servers); err != nil {
		return nil, fmt.Errorf("error decoding get servers response body: %w", err)
	}

	return servers, nil
}

func (c *Client) CreateServer(ctx context.Context, body CreateServerRequestBody) (*Server, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "servers"

	req, err := http.NewRequestWithContext(ctx, "POST", serverURL.String(), bodyReader)
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
			return nil, fmt.Errorf("error decoding create server error response body: %w", err)
		}

		return nil, fmt.Errorf("error creating server: %w", apiError)
	}

	server := Server{}

	if err = json.NewDecoder(res.Body).Decode(&server); err != nil {
		return nil, fmt.Errorf("error decoding create server response body: %w", err)
	}

	server.CallbackID = res.Header.Get("X-Callback-Id")

	return &server, nil
}

func (c *Client) DeleteServer(ctx context.Context, serverSlug string) (string, error) {
	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return "", err
	}

	serverURL.Path += fmt.Sprintf("servers/%s", serverSlug)

	req, err := http.NewRequestWithContext(ctx, "DELETE", serverURL.String(), nil)
	if err != nil {
		return "", err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		apiError := APIError{}

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return "", fmt.Errorf("error decoding delete server error response body: %w", err)
		}

		return "", fmt.Errorf("error deleting server: %w", apiError)
	}

	return res.Header.Get("X-Callback-Id"), nil
}

func (c *Client) GetServerBySlug(ctx context.Context, serverSlug string) (*Server, error) {
	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s", serverSlug)

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
			return nil, fmt.Errorf("error decoding get server by slug error response body: %w", err)
		}

		return nil, fmt.Errorf("error getting server by slug: %w", apiError)
	}

	var server Server

	if err = json.NewDecoder(res.Body).Decode(&server); err != nil {
		return nil, fmt.Errorf("error decoding get server by slug response body: %w", err)
	}

	return &server, nil
}

func (c *Client) PatchServer(ctx context.Context, serverSlug string, body PatchServerRequestBody) (*Server, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s", serverSlug)

	req, err := http.NewRequestWithContext(ctx, "PATCH", serverURL.String(), bodyReader)
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
			return nil, fmt.Errorf("error decoding patch server error response body: %w", err)
		}

		return nil, fmt.Errorf("error patching server: %w", apiError)
	}

	var server Server

	if err = json.NewDecoder(res.Body).Decode(&server); err != nil {
		return nil, fmt.Errorf("error decoding patch server response body: %w", err)
	}

	return &server, nil
}

func (c *Client) ReinstallServer(ctx context.Context, serverSlug string, body ReinstallServerRequestBody) (string, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	bodyReader = bytes.NewReader(buf)

	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return "", err
	}

	serverURL.Path += fmt.Sprintf("servers/%s/actions/reinstall", serverSlug)

	req, err := http.NewRequestWithContext(ctx, "POST", serverURL.String(), bodyReader)
	if err != nil {
		return "", err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		apiError := APIError{}

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return "", fmt.Errorf("error decoding reinstall server error response body: %w", err)
		}

		return "", fmt.Errorf("error reinstalling server: %w", apiError)
	}

	return res.Header.Get("X-Callback-Id"), nil
}

func (c *Client) ResizeServer(ctx context.Context, serverSlug string, body ResizeServerRequestBody) (string, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	bodyReader = bytes.NewReader(buf)

	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return "", err
	}

	serverURL.Path += fmt.Sprintf("servers/%s/actions/resize", serverSlug)

	req, err := http.NewRequestWithContext(ctx, "POST", serverURL.String(), bodyReader)
	if err != nil {
		return "", err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		apiError := APIError{}

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return "", fmt.Errorf("error decoding resize server error response body: %w", err)
		}

		return "", fmt.Errorf("error resizing server: %w", apiError)
	}

	return res.Header.Get("X-Callback-Id"), nil
}

func (c *Client) ResizeDryRun(ctx context.Context, serverSlug string, body ResizeServerRequestBody) (*ServerResize, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s/actions/resize/dryrun", serverSlug)

	req, err := http.NewRequestWithContext(ctx, "POST", serverURL.String(), bodyReader)
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
			return nil, fmt.Errorf("error decoding dry run resize server error response body: %w", err)
		}

		return nil, fmt.Errorf("error dry run resizing server: %w", apiError)
	}

	serverResize := ServerResize{}

	if err = json.NewDecoder(res.Body).Decode(&serverResize); err != nil {
		return nil, fmt.Errorf("error decoding dry run resize server response body: %w", err)
	}

	return &serverResize, nil
}
