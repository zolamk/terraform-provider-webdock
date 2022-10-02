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

func (c *Client) GetPublicKeys(ctx context.Context) (*PublicKeys, error) {
	req, err := newGetPublicKeysRequest(c.Server)
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

	var publicKeys PublicKeys

	if err = decoder.Decode(&publicKeys); err != nil {
		return nil, err
	}

	return &publicKeys, nil
}

func (c *Client) CreatePublicKey(ctx context.Context, body CreatePublicKeyModel) (*PublicKey, error) {
	req, err := newCreatePublicKeyRequest(c.Server, body)
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
		return nil, fmt.Errorf("error creating public key: %s", res.Status)
	}

	var publicKey PublicKey

	decoder := json.NewDecoder(res.Body)

	if err = decoder.Decode(&publicKey); err != nil {
		return nil, err
	}

	return &publicKey, nil
}

func (c *Client) DeletePublicKey(ctx context.Context, id int64) error {
	req, err := newDeletePublicKeyRequest(c.Server, id)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

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

// newGetPublicKeysRequest generates requests for GetPublicKeys
func newGetPublicKeysRequest(server string) (*http.Request, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "account/publicKeys"

	req, err := http.NewRequest("GET", serverURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// newCreatePublicKeyRequest calls the generic PostPublicKeys builder with application/json body
func newCreatePublicKeyRequest(server string, body CreatePublicKeyModel) (*http.Request, error) {
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

	serverURL.Path += "account/publicKeys"

	return http.NewRequest("POST", serverURL.String(), bodyReader)
}

// newDeletePublicKeyRequest generates requests for DeletePublicKey
func newDeletePublicKeyRequest(server string, id int64) (*http.Request, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += fmt.Sprintf("account/publicKeys/%d", id)

	return http.NewRequest("DELETE", serverURL.String(), nil)
}
