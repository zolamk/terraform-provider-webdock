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

func TestGetAccountInformation(t *testing.T) {
	tests := map[string]struct {
		server       *httptest.Server
		wantErr      error
		ctx          context.Context
		wantResponse *api.AccountInformation
	}{
		"when request errors": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      1,
					"message": "unauthorized request",
				})
			})),
			wantErr: fmt.Errorf("error getting account information: %w", api.APIError{ID: 1, Message: "unauthorized request"}),
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
			wantErr: fmt.Errorf("error decoding account information error response body: %w", &json.UnmarshalTypeError{
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
					"accountBalance": 1,
				})
			})),
			wantErr: fmt.Errorf("error decoding account information response body: %w", &json.UnmarshalTypeError{
				Field:  "accountBalance",
				Struct: "AccountInformation",
				Type:   reflect.TypeOf("1"),
				Value:  "number",
				Offset: 19,
			}),
			ctx: context.Background(),
		},
		"when request is successful": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"accountBalance":            "1996 €",
					"accountBalanceRaw":         "199600",
					"accountBalanceRawCurrency": "EUR",
					"companyName":               "Xula",
					"isTeamMember":              true,
					"teamLeader":                "xula@example.com",
					"userAvatar":                "",
					"userEmail":                 "zola@example.com",
					"userId":                    1,
					"userName":                  "Example Name",
				})
			})),
			ctx: context.Background(),
			wantResponse: &api.AccountInformation{
				AccountBalance:            "1996 €",
				AccountBalanceRaw:         "199600",
				AccountBalanceRawCurrency: "EUR",
				CompanyName:               "Xula",
				IsTeamMember:              true,
				TeamLeader:                "xula@example.com",
				UserAvatar:                "",
				UserEmail:                 "zola@example.com",
				UserId:                    json.Number("1"),
				UserName:                  "Example Name",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer test.server.Close()

			client, err := api.NewClient(test.server.URL)

			assert.Nil(t, err)

			account, err := client.GetAccountInformation(test.ctx)

			assert.Equal(t, test.wantErr, err)

			assert.Equal(t, test.wantResponse, account)
		})
	}
}
