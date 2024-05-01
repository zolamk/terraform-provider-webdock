package resource_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zolamk/terraform-provider-webdock/api"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/test/mocks"
	"github.com/zolamk/terraform-provider-webdock/webdock/resource"
)

func TestResourceWebdockPublicKeyCreate(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewClientInterface(t)
	mockErr := errors.New("mock error")
	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"when create public key fails": {
			rd:    resource.PublicKey().Data(&terraform.InstanceState{}),
			diags: diag.FromErr(mockErr),
			mock: func() {
				client.On("CreatePublicKey", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
		},
		"success": {
			rd: resource.PublicKey().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("CreatePublicKey", ctx, mock.Anything).Once().Return(&api.PublicKey{
					Id:      json.Number("1"),
					Created: "2022-05-03T15:05:34+03:00",
					Key:     "test",
					Name:    "test",
				}, nil)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := resource.PublicKey().CreateContext(ctx, test.rd, config.NewCombinedConfig(&config.Config{
				ServerUpPort: 2200,
			}, client))

			assert.Equal(t, test.diags, diags)
		})
	}
}

func TestResourceWebdockPublicKeyDelete(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewClientInterface(t)
	mockErr := errors.New("mock error")
	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"when converting public key id to int64 fails": {
			rd:    resource.PublicKey().Data(&terraform.InstanceState{}),
			diags: diag.Errorf("error converting public key id to int64: strconv.ParseInt: parsing \"\": invalid syntax"),
			mock:  func() {},
		},
		"when delete public key fails": {
			rd: resource.PublicKey().Data(&terraform.InstanceState{
				ID: "1",
			}),
			diags: diag.FromErr(mockErr),
			mock: func() {
				client.On("DeletePublicKey", ctx, mock.Anything).Once().Return(mockErr)
			},
		},
		"success": {
			rd: resource.PublicKey().Data(&terraform.InstanceState{
				ID: "1",
			}),
			mock: func() {
				client.On("DeletePublicKey", ctx, mock.Anything).Once().Return(nil)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := resource.PublicKey().DeleteContext(ctx, test.rd, config.NewCombinedConfig(&config.Config{
				ServerUpPort: 2200,
			}, client))

			assert.Equal(t, test.diags, diags)
		})
	}
}

func TestResourceWebdockPublicKeyRead(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewClientInterface(t)
	mockErr := errors.New("mock error")
	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"when get public keys fails": {
			rd:    resource.PublicKey().Data(&terraform.InstanceState{}),
			diags: diag.Errorf("error getting public key: %v", mockErr),
			mock: func() {
				client.On("GetPublicKeys", ctx).Once().Return(nil, mockErr)
			},
		},
		"when public key is not found": {
			rd: resource.PublicKey().Data(&terraform.InstanceState{
				ID: "3",
			}),
			diags: diag.Errorf("error getting public key: not found"),
			mock: func() {
				client.On("GetPublicKeys", ctx).Once().Return(api.PublicKeys{
					{
						Id:      json.Number("1"),
						Created: "2022-03-20T04:32:12+03:00",
						Key:     "test",
						Name:    "test",
					},
					{
						Id:      json.Number("2"),
						Created: "2022-03-20T04:32:12+03:00",
						Key:     "test2",
						Name:    "test2",
					},
				}, nil)
			},
		},
		"when public keys is nil": {
			rd: resource.PublicKey().Data(&terraform.InstanceState{
				ID: "3",
			}),
			diags: diag.Errorf("error getting public key: not found"),
			mock: func() {
				client.On("GetPublicKeys", ctx).Once().Return(nil, nil)
			},
		},
		"success": {
			rd: resource.PublicKey().Data(&terraform.InstanceState{
				ID: "2",
			}),
			mock: func() {
				client.On("GetPublicKeys", ctx).Once().Return(api.PublicKeys{
					{
						Id:      json.Number("1"),
						Created: "2022-03-20T04:32:12+03:00",
						Key:     "test",
						Name:    "test",
					},
					{
						Id:      json.Number("2"),
						Created: "2022-03-20T04:32:12+03:00",
						Key:     "test2",
						Name:    "test2",
					},
				}, nil)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := resource.PublicKey().ReadContext(ctx, test.rd, config.NewCombinedConfig(&config.Config{
				ServerUpPort: 2200,
			}, client))

			assert.Equal(t, test.diags, diags)
		})
	}
}
