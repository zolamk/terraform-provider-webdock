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

func TestGetEvents(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		params       *api.GetEventsParams
		wantResponse api.Events
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error getting events: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
			ctx:     context.Background(),
		},
		"when error decoding error response": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      true,
					"message": "unexpected error response",
				})
			})),
			wantErr: fmt.Errorf("error decoding get events error response body: %w", &json.UnmarshalTypeError{
				Field:  "id",
				Struct: "APIError",
				Type:   reflect.TypeOf(1),
				Value:  "bool",
				Offset: 10,
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
			wantErr: fmt.Errorf("error decoding get events response body: %w", &json.UnmarshalTypeError{
				Field:  "id",
				Struct: "EventLog",
				Type:   reflect.TypeOf(json.Number("1")),
				Value:  "bool",
				Offset: 11,
			}),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"id":         1,
						"startTime":  "25/10/2022 05:11:34",
						"endTime":    "25/10/2022 05:13:34",
						"callbackId": "3AtrHlEVRg",
						"serverSlug": "server1",
						"eventType":  "start",
						"status":     "waiting",
					},
					{
						"id":         2,
						"startTime":  "25/10/2022 05:11:34",
						"endTime":    "25/10/2022 05:13:34",
						"callbackId": "3AtrHlEVRg",
						"serverSlug": "server1",
						"eventType":  "stop",
						"status":     "waiting",
					},
				})
			})),
			ctx: context.Background(),
			wantResponse: api.Events{
				{
					Id:         json.Number("1"),
					StartTime:  "25/10/2022 05:11:34",
					EndTime:    "25/10/2022 05:13:34",
					CallbackId: "3AtrHlEVRg",
					ServerSlug: "server1",
					EventType:  "start",
					Status:     "waiting",
				},
				{
					Id:         json.Number("2"),
					StartTime:  "25/10/2022 05:11:34",
					EndTime:    "25/10/2022 05:13:34",
					CallbackId: "3AtrHlEVRg",
					ServerSlug: "server1",
					EventType:  "stop",
					Status:     "waiting",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			events, err := client.GetEvents(test.ctx, test.params)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, events)
		})
	}
}
