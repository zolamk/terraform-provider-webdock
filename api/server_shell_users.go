package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ShellUser struct {
	ID         json.Number `json:"id,omitempty"`
	Username   string      `json:"username,omitempty"`
	Group      string      `json:"group,omitempty"`
	Shell      string      `json:"shell,omitempty"`
	PublicKeys []int       `json:"publicKeys,omitempty"`
	Created    time.Time   `json:"created,omitempty"`
	CallbackID string      `json:"-"`
}

type ShellUsers []ShellUser

func (c *Client) GetShellUsers(ctx context.Context, serverSlug string) (ShellUsers, error) {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("servers/%s/shellUsers", serverSlug)

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var shellUsers ShellUsers

	if err = json.NewDecoder(res.Body).Decode(&shellUsers); err != nil {
		return nil, err
	}

	return shellUsers, nil
}

func (c *Client) CreateShellUser(ctx context.Context, serverSlug string, shellUser *ShellUser) (*ShellUser, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(shellUser)
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
		return nil, fmt.Errorf("error creating shell users: %s", res.Status)
	}

	if err = json.NewDecoder(res.Body).Decode(shellUser); err != nil {
		return nil, err
	}

	shellUser.CallbackID = res.Header.Get("X-Callback-ID")

	return shellUser, nil
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
		return "", fmt.Errorf("error deleting shell user: %s", res.Status)
	}

	return res.Header.Get("X-Callback-ID"), nil
}

func (c *Client) UpdateShellUserPublicKeys(ctx context.Context, serverSlug string, shellUserID int64, publicKeys []int) (*ShellUser, error) {
	var bodyReader io.Reader

	shellUser := &ShellUser{
		PublicKeys: publicKeys,
	}

	buf, err := json.Marshal(shellUser)
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
		return nil, fmt.Errorf("error updating shell user: %s", res.Status)
	}

	if err = json.NewDecoder(res.Body).Decode(shellUser); err != nil {
		return nil, err
	}

	shellUser.CallbackID = res.Header.Get("X-Callback-ID")

	return shellUser, nil
}
