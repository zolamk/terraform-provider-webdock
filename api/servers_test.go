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

func TestGetServers(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		wantResponse api.Servers
		serverSlug   string
		params       api.GetServersParams
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error getting servers: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding get servers error response body: %w", &json.UnmarshalTypeError{
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
						"slug": true,
					},
				})
			})),
			ctx: context.Background(),
			wantErr: fmt.Errorf("error decoding get servers response body: %w", &json.UnmarshalTypeError{
				Field:  "slug",
				Struct: "Server",
				Type:   reflect.TypeOf(""),
				Value:  "bool",
				Offset: 13,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"SSHPasswordAuthEnabled": false,
						"WordPressLockDown":      true,
						"aliases":                []string{"server.vps.webdock.cloud"},
						"date":                   "19/09/2022 17:53:14",
						"description":            "Server",
						"image":                  "image_1",
						"ipv4":                   "55.113.30.116",
						"ipv6":                   "5e35:adf4:c51c:d9b3:da85:86f3:ef98:1047",
						"location":               "ad",
						"name":                   "Server",
						"nextActionDate":         "04/10/2022 10:12:28",
						"notes":                  "some note",
						"profile":                "profile_1",
						"slug":                   "server_1",
						"snapshotRunTime":        0,
						"status":                 "started",
						"virtualization":         "container",
						"webServer":              "nginx",
					},
				})
			})),
			ctx: context.Background(),
			wantResponse: api.Servers{
				{
					SSHPasswordAuthEnabled: false,
					WordPressLockDown:      true,
					Aliases:                []string{"server.vps.webdock.cloud"},
					Date:                   "19/09/2022 17:53:14",
					Description:            "Server",
					Image:                  "image_1",
					Ipv4:                   "55.113.30.116",
					Ipv6:                   "5e35:adf4:c51c:d9b3:da85:86f3:ef98:1047",
					Location:               "ad",
					Name:                   "Server",
					NextActionDate:         "04/10/2022 10:12:28",
					Notes:                  "some note",
					Profile:                "profile_1",
					Slug:                   "server_1",
					SnapshotRunTime:        0,
					Status:                 "started",
					Virtualization:         "container",
					WebServer:              "nginx",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			servers, err := client.GetServers(test.ctx, test.params)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, servers)
		})
	}
}

func TestCreateServer(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		req          api.CreateServerRequestBody
		wantResponse *api.Server
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error creating server: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding create server error response body: %w", &json.UnmarshalTypeError{
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
					"slug": true,
				})
			})),
			ctx: context.Background(),
			wantErr: fmt.Errorf("error decoding create server response body: %w", &json.UnmarshalTypeError{
				Field:  "slug",
				Struct: "Server",
				Type:   reflect.TypeOf(""),
				Value:  "bool",
				Offset: 12,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				server := api.CreateServerRequestBody{}

				_ = json.NewDecoder(r.Body).Decode(&server)

				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"image":          server.ImageSlug,
					"location":       server.LocationId,
					"name":           server.Name,
					"profile":        server.ProfileSlug,
					"slug":           server.Slug,
					"virtualization": server.Virtualization,
				})
			})),
			ctx: context.Background(),
			req: api.CreateServerRequestBody{
				ImageSlug:      "image_1",
				LocationId:     "ad",
				Name:           "Server",
				ProfileSlug:    "profile_1",
				Slug:           "server",
				SnapshotId:     10,
				Virtualization: "kvm",
			},
			wantResponse: &api.Server{
				Image:          "image_1",
				Location:       "ad",
				Name:           "Server",
				Profile:        "profile_1",
				Slug:           "server",
				Virtualization: "kvm",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			publicKeys, err := client.CreateServer(test.ctx, test.req)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, publicKeys)
		})
	}
}

func TestDeleteServer(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		serverSlug   string
		wantResponse string
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "server not found",
				})
			})),
			wantErr: fmt.Errorf("error deleting server: %w", api.APIError{ID: 1, Message: "server not found"}),
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
			wantErr: fmt.Errorf("error decoding delete server error response body: %w", &json.UnmarshalTypeError{
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
			wantResponse: "esn0WghLJ3",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			callbackID, err := client.DeleteServer(test.ctx, test.serverSlug)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, callbackID)
		})
	}
}

func TestGetServerBySlug(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		wantResponse *api.Server
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
			wantErr: fmt.Errorf("error getting server by slug: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding get server by slug error response body: %w", &json.UnmarshalTypeError{
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
					"slug": true,
				})
			})),
			ctx: context.Background(),
			wantErr: fmt.Errorf("error decoding get server by slug response body: %w", &json.UnmarshalTypeError{
				Field:  "slug",
				Struct: "Server",
				Type:   reflect.TypeOf(""),
				Value:  "bool",
				Offset: 12,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"SSHPasswordAuthEnabled": false,
					"WordPressLockDown":      true,
					"aliases":                []string{"server.vps.webdock.cloud"},
					"date":                   "19/09/2022 17:53:14",
					"description":            "Server",
					"image":                  "image_1",
					"ipv4":                   "55.113.30.116",
					"ipv6":                   "5e35:adf4:c51c:d9b3:da85:86f3:ef98:1047",
					"location":               "ad",
					"name":                   "Server",
					"nextActionDate":         "04/10/2022 10:12:28",
					"notes":                  "some note",
					"profile":                "profile_1",
					"slug":                   "server_1",
					"snapshotRunTime":        0,
					"status":                 "started",
					"virtualization":         "container",
					"webServer":              "nginx",
				},
				)
			})),
			ctx: context.Background(),
			wantResponse: &api.Server{
				SSHPasswordAuthEnabled: false,
				WordPressLockDown:      true,
				Aliases:                []string{"server.vps.webdock.cloud"},
				Date:                   "19/09/2022 17:53:14",
				Description:            "Server",
				Image:                  "image_1",
				Ipv4:                   "55.113.30.116",
				Ipv6:                   "5e35:adf4:c51c:d9b3:da85:86f3:ef98:1047",
				Location:               "ad",
				Name:                   "Server",
				NextActionDate:         "04/10/2022 10:12:28",
				Notes:                  "some note",
				Profile:                "profile_1",
				Slug:                   "server_1",
				SnapshotRunTime:        0,
				Status:                 "started",
				Virtualization:         "container",
				WebServer:              "nginx",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			server, err := client.GetServerBySlug(test.ctx, test.serverSlug)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, server)
		})
	}
}

func TestPatchServer(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		serverSlug   string
		req          api.PatchServerRequestBody
		wantResponse *api.Server
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error patching server: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding patch server error response body: %w", &json.UnmarshalTypeError{
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
					"slug": true,
				})
			})),
			ctx: context.Background(),
			wantErr: fmt.Errorf("error decoding patch server response body: %w", &json.UnmarshalTypeError{
				Field:  "slug",
				Struct: "Server",
				Type:   reflect.TypeOf(""),
				Value:  "bool",
				Offset: 12,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				server := api.PatchServerRequestBody{}

				_ = json.NewDecoder(r.Body).Decode(&server)

				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"name":           server.Name,
					"description":    server.Description,
					"nextActionDate": server.NextActionDate,
					"notes":          server.Notes,
				})
			})),
			ctx: context.Background(),
			req: api.PatchServerRequestBody{
				Name:           "Server",
				Description:    "server description",
				NextActionDate: "17/08/2022 02:39:22",
				Notes:          "notes",
			},
			wantResponse: &api.Server{
				Name:           "Server",
				Description:    "server description",
				NextActionDate: "17/08/2022 02:39:22",
				Notes:          "notes",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			callbackID, err := client.PatchServer(test.ctx, test.serverSlug, test.req)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, callbackID)
		})
	}
}

func TestReinstallServer(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		serverSlug   string
		req          api.ReinstallServerRequestBody
		wantResponse string
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error reinstalling server: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding reinstall server error response body: %w", &json.UnmarshalTypeError{
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
				server := api.ReinstallServerRequestBody{}

				_ = json.NewDecoder(r.Body).Decode(&server)

				w.Header().Add("X-Callback-ID", "nLVhulfqCy")

				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"image": server.ImageSlug,
				})
			})),
			ctx: context.Background(),
			req: api.ReinstallServerRequestBody{
				ImageSlug: "image_2",
			},
			wantResponse: "nLVhulfqCy",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			callbackID, err := client.ReinstallServer(test.ctx, test.serverSlug, test.req)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, callbackID)
		})
	}
}

func TestResizeServer(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		serverSlug   string
		req          api.ResizeServerRequestBody
		wantResponse string
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error resizing server: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding resize server error response body: %w", &json.UnmarshalTypeError{
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
				server := api.ResizeServerRequestBody{}

				_ = json.NewDecoder(r.Body).Decode(&server)

				w.Header().Add("X-Callback-ID", "nLVhulfqCy")

				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"profile": server.ProfileSlug,
				})
			})),
			ctx: context.Background(),
			req: api.ResizeServerRequestBody{
				ProfileSlug: "profile_2",
			},
			wantResponse: "nLVhulfqCy",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			callbackID, err := client.ResizeServer(test.ctx, test.serverSlug, test.req)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, callbackID)
		})
	}
}

func TestResizeDryRunServer(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		serverSlug   string
		req          api.ResizeServerRequestBody
		wantResponse *api.ServerResize
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error dry run resizing server: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding dry run resize server error response body: %w", &json.UnmarshalTypeError{
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
					"warnings": true,
				})
			})),
			wantErr: fmt.Errorf("error decoding dry run resize server response body: %w", &json.UnmarshalTypeError{
				Field:  "warnings",
				Struct: "ServerResize",
				Type:   reflect.TypeOf([]api.Warning{}),
				Value:  "bool",
				Offset: 16,
			}),
			ctx: context.Background(),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				server := api.ResizeServerRequestBody{}

				_ = json.NewDecoder(r.Body).Decode(&server)

				w.Header().Add("X-Callback-ID", "nLVhulfqCy")

				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"warnings": []map[string]interface{}{
						{
							"data": map[string]interface{}{
								"x": "y",
							},
							"message": "message",
							"type":    "type",
						},
					},
					"chargeSummary": map[string]interface{}{
						"isRefund": true,
						"items": []map[string]interface{}{
							{
								"price": map[string]interface{}{
									"amount":   10,
									"currency": "EUR",
								},
							},
						},
						"total": map[string]interface{}{
							"subTotal": map[string]interface{}{
								"amount":   5,
								"currency": "EUR",
							},
							"total": map[string]interface{}{
								"amount":   5,
								"currency": "EUR",
							},
							"vat": map[string]interface{}{
								"amount":   5,
								"currency": "EUR",
							},
						},
					},
				})
			})),
			ctx: context.Background(),
			req: api.ResizeServerRequestBody{
				ProfileSlug: "profile_2",
			},
			wantResponse: &api.ServerResize{
				Warnings: []api.Warning{
					{
						Data: map[string]interface{}{
							"x": "y",
						},
						Message: "message",
						Type:    "type",
					},
				},
				ChargeSummary: &api.ChargeSummary{
					IsRefund: true,
					Items: []api.ChargeSummaryItem{
						{
							Price: api.Price{
								Amount:   10,
								Currency: "EUR",
							},
						},
					},
					Total: api.ChargeSummaryTotal{
						SubTotal: api.Price{
							Amount:   5,
							Currency: "EUR",
						},
						Total: api.Price{
							Amount:   5,
							Currency: "EUR",
						},
						Vat: api.Price{
							Amount:   5,
							Currency: "EUR",
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			serverResize, err := client.ResizeDryRun(test.ctx, test.serverSlug, test.req)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, serverResize)
		})
	}
}
