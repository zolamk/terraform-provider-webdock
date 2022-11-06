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
type CreatePublicKeyRequestBody struct {
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

	if errorStatus(res.StatusCode) {
		apiError := APIError{}

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error decoding get public keys error response body: %w", err)
		}

		return nil, fmt.Errorf("error getting public keys: %w", apiError)
	}

	var publicKeys PublicKeys

	if err = json.NewDecoder(res.Body).Decode(&publicKeys); err != nil {
		return nil, fmt.Errorf("error decoding get public keys response body: %w", err)
	}

	return publicKeys, nil
}

func (c *Client) CreatePublicKey(ctx context.Context, body CreatePublicKeyRequestBody) (*PublicKey, error) {
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
		apiError := APIError{}

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error decoding create public key error response body: %w", err)
		}

		return nil, fmt.Errorf("error creating public key: %w", apiError)
	}

	publicKey := &PublicKey{}

	if err = json.NewDecoder(res.Body).Decode(publicKey); err != nil {
		return nil, fmt.Errorf("error decoding create public key response body: %w", err)
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
		apiError := APIError{}

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return fmt.Errorf("error decoding delete public key error response body: %w", err)
		}

		return fmt.Errorf("error deleting public key: %w", apiError)
	}

	return nil
}
