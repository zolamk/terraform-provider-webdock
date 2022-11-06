package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zolamk/terraform-provider-webdock/api"
)

func TestGetPublicKeys(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		wantResponse api.PublicKeys
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error getting public keys: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
			ctx:     context.Background(),
		},
		"when error decoding error response": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      "1",
					"message": "unexpected error response",
				})
			})),
			wantErr: fmt.Errorf("error decoding get public keys error response body: %w", &json.UnmarshalTypeError{
				Field:  "id",
				Struct: "APIError",
				Type:   reflect.TypeOf(1),
				Value:  "string",
				Offset: 9,
			}),
			ctx: context.Background(),
		},
		"when error decoding response": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"id": true,
					},
				})
			})),
			ctx: context.Background(),
			wantErr: fmt.Errorf("error decoding get public keys response body: %w", &json.UnmarshalTypeError{
				Field:  "id",
				Struct: "PublicKey",
				Type:   reflect.TypeOf(json.Number("0")),
				Value:  "bool",
				Offset: 11,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"id":      1,
						"name":    "Public Key 1",
						"created": "27/07/2022 11:29:22",
						"key":     "public key 1 content",
					},
					{
						"id":      2,
						"name":    "Public Key 2",
						"created": "27/08/2022 11:29:22",
						"key":     "public key 2 content",
					},
				})
			})),
			ctx: context.Background(),
			wantResponse: api.PublicKeys{
				{
					Id:      json.Number("1"),
					Name:    "Public Key 1",
					Created: "27/07/2022 11:29:22",
					Key:     "public key 1 content",
				},
				{
					Id:      json.Number("2"),
					Name:    "Public Key 2",
					Created: "27/08/2022 11:29:22",
					Key:     "public key 2 content",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			publicKeys, err := client.GetPublicKeys(test.ctx)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, publicKeys)
		})
	}
}

func TestCreatePublicKey(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		req          api.CreatePublicKeyRequestBody
		wantResponse *api.PublicKey
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error creating public key: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
			ctx:     context.Background(),
		},
		"when error decoding error response": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      "1",
					"message": "unexpected error response",
				})
			})),
			wantErr: fmt.Errorf("error decoding create public key error response body: %w", &json.UnmarshalTypeError{
				Field:  "id",
				Struct: "APIError",
				Type:   reflect.TypeOf(1),
				Value:  "string",
				Offset: 9,
			}),
			ctx: context.Background(),
		},
		"when error decoding response": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id": true,
				})
			})),
			ctx: context.Background(),
			wantErr: fmt.Errorf("error decoding create public key response body: %w", &json.UnmarshalTypeError{
				Field:  "id",
				Struct: "PublicKey",
				Type:   reflect.TypeOf(json.Number("0")),
				Value:  "bool",
				Offset: 10,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				publicKey := api.CreatePublicKeyRequestBody{}

				_ = json.NewDecoder(r.Body).Decode(&publicKey)

				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"name":    publicKey.Name,
					"created": "27/07/2022 11:29:22",
					"key":     publicKey.PublicKey,
				})
			})),
			ctx: context.Background(),
			req: api.CreatePublicKeyRequestBody{
				Name:      "Public Key 1",
				PublicKey: "public key 1 content",
			},
			wantResponse: &api.PublicKey{
				Id:      json.Number("1"),
				Name:    "Public Key 1",
				Created: "27/07/2022 11:29:22",
				Key:     "public key 1 content",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			publicKeys, err := client.CreatePublicKey(test.ctx, test.req)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, publicKeys)
		})
	}
}

func TestDeletePublicKey(t *testing.T) {
	tests := map[string]struct {
		server  *httptest.Server
		wantErr error
		ctx     context.Context
		id      int64
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "public key not found",
				})
			})),
			wantErr: fmt.Errorf("error deleting public key: %w", api.APIError{ID: 1, Message: "public key not found"}),
			ctx:     context.Background(),
		},
		"when error decoding error response": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      "1",
					"message": "unexpected error response",
				})
			})),
			wantErr: fmt.Errorf("error decoding delete public key error response body: %w", &json.UnmarshalTypeError{
				Field:  "id",
				Struct: "APIError",
				Type:   reflect.TypeOf(1),
				Value:  "string",
				Offset: 9,
			}),
			ctx: context.Background(),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})),
			ctx: context.Background(),
			id:  1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			err = client.DeletePublicKey(test.ctx, test.id)

			assert.Equal(t, test.wantErr, err)
		})
	}
}
