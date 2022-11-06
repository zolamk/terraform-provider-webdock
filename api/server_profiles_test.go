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

func TestGetServerProfiles(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		wantResponse api.ServerProfiles
		params       api.GetServersProfilesParams
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error getting server profiles: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding get server profiles error response body: %w", &json.UnmarshalTypeError{
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
			wantErr: fmt.Errorf("error decoding get server profiles response body: %w", &json.UnmarshalTypeError{
				Field:  "name",
				Struct: "ServerProfile",
				Type:   reflect.TypeOf(""),
				Value:  "bool",
				Offset: 13,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"cpu": map[string]interface{}{
							"cores":   1,
							"threads": 2,
						},
						"disk": 23842,
						"name": "SSD Nano4",
						"ram":  1907,
						"slug": "webdocknano4-2022",
					},
				})
			})),
			ctx: context.Background(),
			wantResponse: api.ServerProfiles{
				{
					CPU: api.CPU{
						Cores:   1,
						Threads: 2,
					},
					Disk: 23842,
					Name: "SSD Nano4",
					RAM:  1907,
					Slug: "webdocknano4-2022",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			serverProfiles, err := client.GetServersProfiles(test.ctx, test.params)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, serverProfiles)
		})
	}
}
