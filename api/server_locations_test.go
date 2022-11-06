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

func TestGetServerLocations(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		wantResponse api.ServerLocations
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error getting server locations: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding get server locations error response body: %w", &json.UnmarshalTypeError{
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
			wantErr: fmt.Errorf("error decoding get server locations response body: %w", &json.UnmarshalTypeError{
				Field:  "name",
				Struct: "ServerLocation",
				Type:   reflect.TypeOf(""),
				Value:  "bool",
				Offset: 13,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"city":        "Addis Ababa",
						"country":     "Ethiopia",
						"description": "Server will be located in our datacenter in Addis Ababa, Ethiopia",
						"icon":        "icon",
						"id":          "ad",
						"name":        "Africa",
					},
					{
						"city":        "Montreal",
						"country":     "Canada",
						"description": "Server will be located in our datacenter in Montreal, Canada",
						"icon":        "https://api.webdock.io/concrete/images/countries/ca.png",
						"id":          "ca",
						"name":        "North America",
					},
					{
						"city":        "Helsinki",
						"country":     "Finland",
						"description": "Server will be located in our datacenter in Helsinki, Finland",
						"icon":        "https://api.webdock.io/concrete/images/countries/europeanunion.png",
						"id":          "fi",
						"name":        "Europe",
					},
				})
			})),
			ctx: context.Background(),
			wantResponse: api.ServerLocations{
				{
					ID:          "ad",
					City:        "Addis Ababa",
					Country:     "Ethiopia",
					Description: "Server will be located in our datacenter in Addis Ababa, Ethiopia",
					Icon:        "icon",
					Name:        "Africa",
				}, {
					City:        "Montreal",
					Country:     "Canada",
					Description: "Server will be located in our datacenter in Montreal, Canada",
					Icon:        "https://api.webdock.io/concrete/images/countries/ca.png",
					ID:          "ca",
					Name:        "North America",
				},
				{
					City:        "Helsinki",
					Country:     "Finland",
					Description: "Server will be located in our datacenter in Helsinki, Finland",
					Icon:        "https://api.webdock.io/concrete/images/countries/europeanunion.png",
					ID:          "fi",
					Name:        "Europe",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			serverLocations, err := client.GetServersLocations(test.ctx)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, serverLocations)
		})
	}
}
