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

// Defines values for ServerDTOStatus.
const (
	ServerDTOStatusError        ServerDTOStatus = "error"
	ServerDTOStatusProvisioning ServerDTOStatus = "provisioning"
	ServerDTOStatusRebooting    ServerDTOStatus = "rebooting"
	ServerDTOStatusReinstalling ServerDTOStatus = "reinstalling"
	ServerDTOStatusRunning      ServerDTOStatus = "running"
	ServerDTOStatusStarting     ServerDTOStatus = "starting"
	ServerDTOStatusStopped      ServerDTOStatus = "stopped"
	ServerDTOStatusStopping     ServerDTOStatus = "stopping"
)

// Charge summary items model
type ChargeSummaryItem struct {
	// Price model
	Price *PriceDTO `json:"price,omitempty"`

	// Charge text
	Text *string `json:"text,omitempty"`
}

// Server resize model
type ChargeSummary struct {
	// True if the amount will be refunded
	IsRefund *bool `json:"isRefund,omitempty"`

	// List of charges to be applied or refunded
	Items *[]ChargeSummaryItem `json:"items,omitempty"`

	// Server resize model
	Total *ChargeSummaryTotalDTO `json:"total,omitempty"`
}

// PatchServer model
type PatchServerModelDTO struct {
	// Description of the server
	Description string `json:"description"`

	// Name of the server
	Name string `json:"name"`

	// Next action of the server
	NextActionDate *string `json:"nextActionDate"`

	// Internal notes or comments regarding the server
	Notes string `json:"notes"`
}

// Reinstall Server model
type ReinstallServerModelDTO struct {
	// Image slug of the image you want to reload the server with. Any image listed for the server location from /images is valid.
	ImageSlug string `json:"imageSlug"`
}

// Resize server model
type ResizeServerModelDTO struct {
	// Profile slug to resize to
	ProfileSlug string `json:"profileSlug"`
}

// PatchServer model
type PatchServerModel = PatchServerModelDTO

// Post Server model
type PostServerModel = PostServerModelDTO

// Reinstall Server model
type ReinstallServerModel = ReinstallServerModelDTO

// GetServersParams defines parameters for GetServers.
type GetServersParams struct {
	// Filter by current status of the server
	Status string `form:"status,omitempty" json:"status,omitempty"`
}

// CreateServerJSONRequestBody defines body for CreateServer for application/json ContentType.
type CreateServerJSONRequestBody = PostServerModel

// PatchServerJSONRequestBody defines body for PatchServer for application/json ContentType.
type PatchServerJSONRequestBody = PatchServerModel

// ReinstallServerJSONRequestBody defines body for ReinstallServer for application/json ContentType.
type ReinstallServerJSONRequestBody = ReinstallServerModel

// ResizeServerJSONRequestBody defines body for ResizeServer for application/json ContentType.
type ResizeServerJSONRequestBody = ResizeServerModelDTO

// ResizeDryRunJSONRequestBody defines body for ResizeDryRun for application/json ContentType.
type ResizeDryRunJSONRequestBody = ResizeServerModelDTO

// Warning model
type WarningDTO struct {
	// Warning message
	Data *WarningDTO_Data `json:"data,omitempty"`

	// Warning message
	Message *string `json:"message,omitempty"`

	// Warning type
	Type *string `json:"type,omitempty"`
}

// Warning message
type WarningDTO_Data struct {
	AdditionalProperties map[string]interface{} `json:"-"`
}

// Server model
type Server struct {
	// SSH Password Authentication Enabled for this Server
	SSHPasswordAuthEnabled *bool `json:"SSHPasswordAuthEnabled,omitempty" mapstructure:"ssh_password_auth_enabled"`

	// Wordpress lockdown status
	WordPressLockDown *bool `json:"WordPressLockDown,omitempty" mapstructure:"wordpress_lockdown"`

	// Aliases - Domain names for the server as known by Webdock. First entry should be treated as the &quot;Main Domain&quot; for the server.
	Aliases *[]string `json:"aliases,omitempty" mapstructure:"aliases"`

	// Creation date/time
	Date string `json:"date,omitempty" mapstructure:"created_at"`

	// Server Description (what's installed here?) as entered by admin in Server Metadata
	Description *string `json:"description,omitempty" mapstructure:"description"`

	// Server image
	Image string `json:"image,omitempty" mapstructure:"image_slug"`

	// IPv4 address
	Ipv4 string `json:"ipv4" mapstructure:"ipv4"`

	// IPv6 address
	Ipv6 string `json:"ipv6" mapstructure:"ipv6"`

	// Location ID of the server
	Location *string `json:"location" mapstructure:"location_id"`

	// Server name
	Name string `json:"name,omitempty" mapstructure:"name"`

	// Next Action date/time as entered by admin in Server Metadata
	NextActionDate *string `json:"nextActionDate,omitempty" mapstructure:"next_action_date"`

	// Notes as entered by admin in Server Metadata
	Notes *string `json:"notes,omitempty" mapstructure:"notes"`

	// Server profile
	Profile *string `json:"profile" mapstructure:"profile_slug"`

	// Server slug
	Slug string `json:"slug,omitempty" mapstructure:"slug"`

	// Last known snapshot runtime (seconds)
	SnapshotRunTime *int64 `json:"snapshotRunTime,omitempty" mapstructure:"snapshot_run_time"`

	// Server status
	Status *ServerDTOStatus `json:"status,omitempty" mapstructure:"status"`

	// Webserver type
	WebServer *string `json:"webServer,omitempty" mapstructure:"webserver"`

	CallbackID string `json:"-" mapstructure:"-"`
}

// Servers is a collection of Server
type Servers []Server

// Server status
type ServerDTOStatus string

// Server resize model
type ServerResize struct {
	ChargeSummary *ChargeSummary `json:"chargeSummary,omitempty"`
	Warnings      *[]WarningDTO  `json:"warnings,omitempty"`
}

// Server resize model
type ChargeSummaryTotalDTO struct {
	// Price model
	SubTotal *PriceDTO `json:"subTotal,omitempty"`

	// Price model
	Total *PriceDTO `json:"total,omitempty"`

	// Price model
	Vat *PriceDTO `json:"vat,omitempty"`
}

// Post Server model
type PostServerModelDTO struct {
	// Slug of the server image. Get this from the /images endpoint. You must pass either this parameter or snapshotId
	ImageSlug string `json:"imageSlug,omitempty"`

	// ID of the location. Get this from the /locations endpoint.
	LocationId string `json:"locationId"`

	// Name of the server
	Name string `json:"name"`

	// Slug of the server profile. Get this from the /profiles endpoint.
	ProfileSlug string `json:"profileSlug"`

	// Suggested Slug (shortname) of the server. Up to 12 alphanumeric chars. This slug is effectively your server ID and anything submitted in this field is merely your suggestion for a slug. If omitted or the suggested slug is already taken, the system will automatically generate an appropriate unique slug based on your server Name or suggestion. Always check the return from this method to determine the actual slug your server ended up receiving.
	Slug *string `json:"slug,omitempty"`

	// SnapshotID from which to create the server. Get this from the /servers/{serverSlug}/snapshots endpoint. You must pass either this parameter or imageSlug.
	SnapshotId *int64 `json:"snapshotId,omitempty"`
}

func (c *Client) GetServers(ctx context.Context, params *GetServersParams, reqEditors ...RequestEditorFn) (*Servers, error) {
	req, err := newGetServersRequest(c.Server, params)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		return nil, fmt.Errorf("error getting server: %s", res.Status)
	}

	var servers Servers

	decoder := json.NewDecoder(res.Body)

	if err = decoder.Decode(&servers); err != nil {
		return nil, err
	}

	return &servers, nil
}

func (c *Client) CreateServer(ctx context.Context, body CreateServerJSONRequestBody) (*Server, *string, error) {
	req, err := newCreateServerRequest(c.Server, body)
	if err != nil {
		return nil, nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		return nil, nil, fmt.Errorf("error creating server: %s", res.Status)
	}

	var server Server

	decoder := json.NewDecoder(res.Body)

	if err = decoder.Decode(&server); err != nil {

		return nil, nil, err

	}

	callbackID := res.Header.Get("X-Callback-Id")

	return &server, &callbackID, nil
}

func (c *Client) DeleteServer(ctx context.Context, serverSlug string, reqEditors ...RequestEditorFn) (*string, error) {
	req, err := newDeleteServerRequest(c.Server, serverSlug)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		return nil, fmt.Errorf("error deleting server: %s", res.Status)
	}

	callbackID := res.Header.Get("X-Callback-Id")

	return &callbackID, nil
}

func (c *Client) GetServerBySlug(ctx context.Context, serverSlug string) (*Server, error) {
	req, err := newGetServerBySlugRequest(c.Server, serverSlug)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		return nil, fmt.Errorf("error getting server: %s", res.Status)
	}

	var server Server

	decoder := json.NewDecoder(res.Body)

	if err = decoder.Decode(&server); err != nil {
		return nil, err
	}

	return &server, nil
}

func (c *Client) PatchServer(ctx context.Context, serverSlug string, body PatchServerJSONRequestBody, reqEditors ...RequestEditorFn) (*Server, error) {
	req, err := newPatchServerRequest(c.Server, serverSlug, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		return nil, fmt.Errorf("error patching server: %s", res.Status)
	}

	var server Server

	decoder := json.NewDecoder(res.Body)

	if err = decoder.Decode(&server); err != nil {
		return nil, err
	}

	return &server, nil
}

func (c *Client) ReinstallServer(ctx context.Context, serverSlug string, body ReinstallServerJSONRequestBody) (*string, error) {
	req, err := newReinstallServerRequest(c.Server, serverSlug, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		return nil, fmt.Errorf("error reinstalling server: %s", res.Status)
	}

	callbackID := res.Header.Get("X-Callback-Id")

	return &callbackID, nil
}

func (c *Client) ResizeServer(ctx context.Context, serverSlug string, body ResizeServerJSONRequestBody, reqEditors ...RequestEditorFn) (*string, error) {
	req, err := newResizeServerRequest(c.Server, serverSlug, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		return nil, fmt.Errorf("error resizing server: %s", res.Status)
	}

	callbackID := res.Header.Get("X-Callback-Id")

	return &callbackID, nil
}

func (c *Client) ResizeDryRun(ctx context.Context, serverSlug string, body ResizeDryRunJSONRequestBody, reqEditors ...RequestEditorFn) (*ServerResize, error) {
	req, err := newResizeDryRunRequest(c.Server, serverSlug, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		return nil, fmt.Errorf("error resizing server: %s", res.Status)
	}

	var serverResize ServerResize

	decoder := json.NewDecoder(res.Body)

	if err = decoder.Decode(&serverResize); err != nil {
		return nil, err
	}

	return &serverResize, nil
}

// newGetServersRequest generates requests for GetServers
func newGetServersRequest(server string, params *GetServersParams) (*http.Request, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "servers"

	queryValues, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	serverURL.RawQuery = queryValues.Encode()

	return http.NewRequest("GET", serverURL.String(), nil)
}

// newCreateServerRequest calls the generic CreateServer builder with application/json body
func newCreateServerRequest(server string, body CreateServerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "servers"

	return http.NewRequest("POST", serverURL.String(), bodyReader)
}

// newDeleteServerRequest generates requests for DeleteServer
func newDeleteServerRequest(server string, serverSlug string) (*http.Request, error) {
	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(server)

	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s", serverSlug)

	return http.NewRequest("DELETE", serverURL.String(), nil)
}

// newGetServerBySlugRequest generates requests for GetServerBySlug
func newGetServerBySlugRequest(server string, serverSlug string) (*http.Request, error) {
	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s", serverSlug)

	return http.NewRequest("GET", serverURL.String(), nil)
}

// newPatchServerRequest calls the generic PatchServer builder with application/json body
func newPatchServerRequest(server string, serverSlug string, body PatchServerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s", serverSlug)

	return http.NewRequest("PATCH", serverURL.String(), bodyReader)
}

// newReinstallServerRequest calls the generic ReinstallServer builder with application/json body
func newReinstallServerRequest(server string, serverSlug string, body ReinstallServerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s/actions/reinstall", serverSlug)

	return http.NewRequest("POST", serverURL.String(), bodyReader)
}

// newResizeServerRequest calls the generic ResizeServer builder with application/json body
func newResizeServerRequest(server string, serverSlug string, body ResizeServerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s/actions/resize", serverSlug)

	return http.NewRequest("POST", serverURL.String(), bodyReader)
}

// newResizeDryRunRequest calls the generic ResizeDryRun builder with application/json body
func newResizeDryRunRequest(server string, serverSlug string, body ResizeDryRunJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	serverSlug = url.PathEscape(serverSlug)

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s/actions/resize/dryrun", serverSlug)

	return http.NewRequest("POST", serverURL.String(), bodyReader)
}
