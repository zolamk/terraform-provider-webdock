package api

import (
	"context"
	"net/http"
	"strings"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientTransport is a custom transport that runs request editor functions before sending the request.
type ClientTransport struct {
	client *Client
}

func (c *ClientTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, r := range c.client.RequestEditors {
		if err := r(req.Context(), req); err != nil {
			return nil, err
		}
	}

	req.Header.Add("Content-Type", "application/json")

	return http.DefaultTransport.RoundTrip(req)
}

type APIError struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}

	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}

	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}

	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{
			Transport: &ClientTransport{&client},
		}
	}

	return &client, nil

}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetPublicKeys request
	GetPublicKeys(ctx context.Context) (PublicKeys, error)

	CreatePublicKey(ctx context.Context, body CreatePublicKeyRequestBody) (*PublicKey, error)

	// DeletePublicKey request
	DeletePublicKey(ctx context.Context, id int64) error

	// GetEvents request
	GetEvents(ctx context.Context, params GetEventsParams) (Events, error)

	// GetServersImages request
	GetServersImages(ctx context.Context) (ServerImages, error)

	// GetServersLocations request
	GetServersLocations(ctx context.Context) (ServerLocations, error)

	// GetServersProfiles request
	GetServersProfiles(ctx context.Context, params GetServersProfilesParams) (ServerProfiles, error)

	// GetServers request
	GetServers(ctx context.Context, params GetServersParams) (Servers, error)

	// CreateServer request
	CreateServer(ctx context.Context, body CreateServerRequestBody) (*Server, error)

	// DeleteServer request
	DeleteServer(ctx context.Context, serverSlug string) (string, error)

	// GetServerBySlug request
	GetServerBySlug(ctx context.Context, serverSlug string) (*Server, error)

	PatchServer(ctx context.Context, serverSlug string, body PatchServerRequestBody) (*Server, error)

	ReinstallServer(ctx context.Context, serverSlug string, body ReinstallServerRequestBody) (string, error)

	ResizeServer(ctx context.Context, serverSlug string, body ResizeServerRequestBody) (string, error)

	ResizeDryRun(ctx context.Context, serverSlug string, body ResizeServerRequestBody) (*ServerResize, error)

	GetShellUsers(ctx context.Context, serverSlug string) (ShellUsers, error)

	CreateShellUser(ctx context.Context, serverSlug string, shellUser CreateShellUserRequestBody) (*ShellUser, error)

	DeleteShellUser(ctx context.Context, serverSlug string, shellUserID int64) (string, error)

	UpdateShellUserPublicKeys(ctx context.Context, serverSlug string, shellUserID int64, publicKeys []int) (*ShellUser, error)
}
