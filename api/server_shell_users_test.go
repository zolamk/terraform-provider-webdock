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

func TestGetShellUsers(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		wantResponse api.ShellUsers
		serverSlug   string
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error getting server shell users: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding get server shell users error response body: %w", &json.UnmarshalTypeError{
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
			wantErr: fmt.Errorf("error decoding get server shell users response body: %w", &json.UnmarshalTypeError{
				Field:  "id",
				Struct: "ShellUser",
				Type:   reflect.TypeOf(json.Number("0")),
				Value:  "bool",
				Offset: 11,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"id":         1,
						"username":   "xula",
						"group":      "sudo",
						"shell":      "/bin/sh",
						"publicKeys": []map[string]interface{}{},
						"created":    "27/07/2022 11:29:22",
					},
					{
						"id":       2,
						"username": "zola",
						"group":    "sudo",
						"shell":    "/bin/bash",
						"publicKeys": []map[string]interface{}{
							{
								"id":   1,
								"key":  "public key 1 content",
								"name": "public key 1",
							},
						},
						"created": "27/08/2022 11:29:22",
					},
				})
			})),
			ctx: context.Background(),
			wantResponse: api.ShellUsers{
				{
					ID:         json.Number("1"),
					Username:   "xula",
					Group:      "sudo",
					Shell:      "/bin/sh",
					PublicKeys: api.PublicKeys{},
					Created:    "27/07/2022 11:29:22",
				},
				{
					ID:       json.Number("2"),
					Username: "zola",
					Group:    "sudo",
					Shell:    "/bin/bash",
					PublicKeys: api.PublicKeys{
						{
							Id:   "1",
							Key:  "public key 1 content",
							Name: "public key 1",
						},
					},
					Created: "27/08/2022 11:29:22",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			shellUsers, err := client.GetShellUsers(test.ctx, test.serverSlug)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, shellUsers)
		})
	}
}

func TestCreateShellUser(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		req          api.CreateShellUserRequestBody
		serverSlug   string
		wantResponse *api.ShellUser
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error creating shell user: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding create shell user error response body: %w", &json.UnmarshalTypeError{
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
			wantErr: fmt.Errorf("error decoding create shell user response body: %w", &json.UnmarshalTypeError{
				Field:  "id",
				Struct: "ShellUser",
				Type:   reflect.TypeOf(json.Number("0")),
				Value:  "bool",
				Offset: 10,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				shellUser := api.CreateShellUserRequestBody{}

				_ = json.NewDecoder(r.Body).Decode(&shellUser)

				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"username": "xula",
					"group":    "sudo",
					"shell":    "/bin/bash",
					"publicKeys": []map[string]interface{}{
						{
							"id":   shellUser.PublicKeys[0],
							"key":  "public key 10 content",
							"name": "public key 10",
						},
					},
					"created": "27/07/2022 11:29:22",
				})
			})),
			ctx: context.Background(),
			req: api.CreateShellUserRequestBody{
				Username:   "xula",
				Password:   "password",
				Group:      "sudo",
				Shell:      "/bin/bash",
				PublicKeys: []int{10},
			},
			wantResponse: &api.ShellUser{
				Username: "xula",
				Group:    "sudo",
				Shell:    "/bin/bash",
				PublicKeys: api.PublicKeys{
					{
						Id:   json.Number("10"),
						Key:  "public key 10 content",
						Name: "public key 10",
					},
				},
				Created: "27/07/2022 11:29:22",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			publicKeys, err := client.CreateShellUser(test.ctx, test.serverSlug, test.req)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, publicKeys)
		})
	}
}

func TestDeleteShellUser(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		id           int64
		serverSlug   string
		wantResponse string
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "server shell user not found",
				})
			})),
			wantErr: fmt.Errorf("error deleting server shell user: %w", api.APIError{ID: 1, Message: "server shell user not found"}),
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
			wantErr: fmt.Errorf("error decoding delete shell user error response body: %w", &json.UnmarshalTypeError{
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
				w.Header().Add("X-Callback-ID", "esn0WghLJ3")
				w.WriteHeader(http.StatusNoContent)
			})),
			ctx:          context.Background(),
			id:           1,
			wantResponse: "esn0WghLJ3",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			callbackID, err := client.DeleteShellUser(test.ctx, test.serverSlug, test.id)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, callbackID)
		})
	}
}

func TestUpdateShellUserPublicKeys(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		serverSlug   string
		shellUserID  int64
		publicKeys   []int
		wantResponse *api.ShellUser
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error updating shell user: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding update shell user error response body: %w", &json.UnmarshalTypeError{
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
			wantErr: fmt.Errorf("error decoding update shell user response body: %w", &json.UnmarshalTypeError{
				Field:  "id",
				Struct: "ShellUser",
				Type:   reflect.TypeOf(json.Number("0")),
				Value:  "bool",
				Offset: 10,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				shellUser := api.CreateShellUserRequestBody{}

				_ = json.NewDecoder(r.Body).Decode(&shellUser)

				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"username": "xula",
					"group":    "sudo",
					"shell":    "/bin/bash",
					"publicKeys": []map[string]interface{}{
						{
							"id":   shellUser.PublicKeys[0],
							"key":  "public key 15 content",
							"name": "public key 15",
						},
					},
					"created": "27/07/2022 11:29:22",
				})
			})),
			ctx:        context.Background(),
			publicKeys: []int64{15},
			wantResponse: &api.ShellUser{
				Username: "xula",
				Group:    "sudo",
				Shell:    "/bin/bash",
				PublicKeys: api.PublicKeys{
					{
						Id:   json.Number("15"),
						Key:  "public key 15 content",
						Name: "public key 15",
					},
				},
				Created: "27/07/2022 11:29:22",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			publicKeys, err := client.UpdateShellUserPublicKeys(test.ctx, test.serverSlug, test.shellUserID, test.publicKeys)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, publicKeys)
		})
	}
}
