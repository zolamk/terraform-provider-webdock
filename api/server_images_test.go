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

func TestGetServerImages(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		wantResponse api.ServerImages
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error getting server images: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding get server images error response body: %w", &json.UnmarshalTypeError{
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
						"name": true,
					},
				})
			})),
			ctx: context.Background(),
			wantErr: fmt.Errorf("error decoding get server images response body: %w", &json.UnmarshalTypeError{
				Field:  "name",
				Struct: "ServerImage",
				Type:   reflect.TypeOf(""),
				Value:  "bool",
				Offset: 13,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"name":       "Image 1",
						"slug":       "image_1",
						"phpVersion": "8.0",
						"webServer":  "nginx",
					},
					{
						"name":      "Image 2",
						"slug":      "image_2",
						"webServer": "apache",
					},
				})
			})),
			ctx: context.Background(),
			wantResponse: api.ServerImages{
				{
					Name:       "Image 1",
					Slug:       "image_1",
					PhpVersion: "8.0",
					WebServer:  "nginx",
				},
				{
					Name:      "Image 2",
					Slug:      "image_2",
					WebServer: "apache",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			serverImages, err := client.GetServersImages(test.ctx)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, serverImages)
		})
	}
}
