package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	errGettingServerShellUsers = "error getting server shell users"
)

type CreateShellUserRequestBody struct {
	Username   string  `json:"username,omitempty"`
	Password   string  `json:"password,omitempty"`
	Group      string  `json:"group,omitempty"`
	Shell      string  `json:"shell,omitempty"`
	PublicKeys []int64 `json:"publicKeys,omitempty"`
}

type ShellUser struct {
	ID         json.Number `json:"id,omitempty" mapstructure:"id"`
	Username   string      `json:"username,omitempty" mapstructure:"username"`
	Password   string      `json:"password,omitempty" mapstructure:"password"`
	Group      string      `json:"group,omitempty" mapstructure:"group"`
	Shell      string      `json:"shell,omitempty" mapstructure:"shell"`
	PublicKeys PublicKeys  `json:"publicKeys,omitempty" mapstructure:"public_keys"`
	Created    string      `json:"created,omitempty" mapstructure:"created_at"`
	CallbackID string      `json:"-" mapstructure:"-"`
}

type ShellUsers []ShellUser

func (c *Client) GetShellUsers(ctx context.Context, serverSlug string) (ShellUsers, error) {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errGettingServerShellUsers, err)
	}

	serverURL.Path += fmt.Sprintf("servers/%s/shellUsers", serverSlug)

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errGettingServerShellUsers, err)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errGettingServerShellUsers, err)
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		apiError := APIError{}

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error decoding get server shell users error response body: %w", err)
		}

		return nil, fmt.Errorf("%s: %w", errGettingServerShellUsers, apiError)
	}

	shellUsers := ShellUsers{}

	if err = json.NewDecoder(res.Body).Decode(&shellUsers); err != nil {
		return nil, fmt.Errorf("error decoding get server shell users response body: %w", err)
	}

	return shellUsers, nil
}

func (c *Client) CreateShellUser(ctx context.Context, serverSlug string, createShellUserBody CreateShellUserRequestBody) (*ShellUser, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(createShellUserBody)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s/shellUsers", serverSlug)

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
			return nil, fmt.Errorf("error decoding create shell user error response body: %w", err)
		}

		return nil, fmt.Errorf("error creating shell user: %w", apiError)
	}

	shellUser := ShellUser{}

	if err = json.NewDecoder(res.Body).Decode(&shellUser); err != nil {
		return nil, fmt.Errorf("error decoding create shell user response body: %w", err)
	}

	shellUser.CallbackID = res.Header.Get("X-Callback-ID")

	return &shellUser, nil
}

func (c *Client) DeleteShellUser(ctx context.Context, serverSlug string, shellUserID int64) (string, error) {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return "", err
	}

	serverURL.Path += fmt.Sprintf("servers/%s/shellUsers/%d", serverSlug, shellUserID)

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
			return "", fmt.Errorf("error decoding delete shell user error response body: %w", err)
		}

		return "", fmt.Errorf("error deleting server shell user: %w", apiError)
	}

	return res.Header.Get("X-Callback-ID"), nil
}

func (c *Client) UpdateShellUserPublicKeys(ctx context.Context, serverSlug string, shellUserID int64, publicKeys []int64) (*ShellUser, error) {
	var bodyReader io.Reader

	shellUserBody := &CreateShellUserRequestBody{
		PublicKeys: publicKeys,
	}

	buf, err := json.Marshal(shellUserBody)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)

	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s/shellUsers/%d", serverSlug, shellUserID)

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
			return nil, fmt.Errorf("error decoding update shell user error response body: %w", err)
		}

		return nil, fmt.Errorf("error updating shell user: %w", apiError)
	}

	shellUser := ShellUser{}

	if err = json.NewDecoder(res.Body).Decode(&shellUser); err != nil {
		return nil, fmt.Errorf("error decoding update shell user response body: %w", err)
	}

	shellUser.CallbackID = res.Header.Get("X-Callback-ID")

	return &shellUser, nil
}
