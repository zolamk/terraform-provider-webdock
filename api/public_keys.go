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

// PublicKey model
type CreatePublicKeyModel struct {
	// PublicKey name
	Name string `json:"name"`

	// PublicKey
	PublicKey string `json:"publicKey"`
}

// PublicKey model
type PublicKey struct {
	// PublicKey ID
	Id json.Number `json:"id,omitempty" mapstructure:"id"`

	// PublicKey creation datetime
	Created string `json:"created,omitempty" mapstructure:"created_at"`

	// PublicKey content
	Key string `json:"key,omitempty" mapstructure:"key"`

	// PublicKey name
	Name string `json:"name,omitempty" mapstructure:"name"`
}

// PublicKeys is a collection of PublicKey
type PublicKeys []PublicKey

func (c *Client) GetPublicKeys(ctx context.Context) (PublicKeys, error) {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "account/publicKeys"

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var publicKeys PublicKeys

	if err = json.NewDecoder(res.Body).Decode(&publicKeys); err != nil {
		return nil, err
	}

	return publicKeys, nil
}

func (c *Client) CreatePublicKey(ctx context.Context, body CreatePublicKeyModel) (*PublicKey, error) {
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

	serverURL.Path += "account/publicKeys"

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
		return nil, fmt.Errorf("error creating public key: %s", res.Status)
	}

	publicKey := &PublicKey{}

	if err = json.NewDecoder(res.Body).Decode(publicKey); err != nil {
		return nil, err
	}

	return publicKey, nil
}

func (c *Client) DeletePublicKey(ctx context.Context, id int64) error {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return err
	}

	serverURL.Path += fmt.Sprintf("account/publicKeys/%d", id)

	req, err := http.NewRequestWithContext(ctx, "DELETE", serverURL.String(), nil)
	if err != nil {
		return err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if errorStatus(res.StatusCode) {
		return fmt.Errorf("error deleting public key: %s", res.Status)
	}

	return nil
}
