package resource_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zolamk/terraform-provider-webdock/api"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/test/mocks"
	"github.com/zolamk/terraform-provider-webdock/webdock/resource"
)

func TestResourceWebdockServerCreate(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewClientInterface(t)
	mockErr := errors.New("mock error")
	l, err := net.Listen("tcp", "127.0.0.1:2200")
	require.Nil(t, err)
	defer l.Close()

	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"when create server fails": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			diags: diag.FromErr(mockErr),
			mock: func() {
				client.On("CreateServer", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
		},
		"when create server fails with too many server error": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			diags: nil,
			mock: func() {
				client.On("CreateServer", ctx, mock.Anything).Twice().Return(nil, &api.APIError{
					ID:      0,
					Message: "You are creating too many servers in too short of a timespan. Please wait a while and try again a bit later.",
				})

				client.On("CreateServer", ctx, mock.Anything).Once().Return(&api.Server{
					SSHPasswordAuthEnabled: true,
					WordPressLockDown:      true,
					Aliases:                []string{"test"},
					Date:                   "2022-12-22T03:54:56+03:00",
					Image:                  "test",
					Ipv4:                   "127.0.0.1",
					Ipv6:                   "8b34:f82b:999a:1ab5:0cad:f252:af94:bf80",
					Location:               "test",
					Name:                   "test",
					Profile:                "test",
					Slug:                   "test",
					SnapshotRunTime:        0,
					Status:                 "provisioning",
					Virtualization:         "containerd",
					WebServer:              "nginx",
					CallbackID:             "callback",
				}, nil)

				client.On("GetEvents", ctx, mock.Anything).Once().Return(api.Events{
					{
						Status: "finished",
					},
				}, nil)
			},
		},
		"when wait for action fails": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			diags: diag.Errorf("server (test) create event (callback) errored: %v", mockErr),
			mock: func() {
				client.On("CreateServer", ctx, mock.Anything).Once().Return(&api.Server{
					SSHPasswordAuthEnabled: true,
					WordPressLockDown:      true,
					Aliases:                []string{"test"},
					Date:                   "2022-12-22T03:54:56+03:00",
					Image:                  "test",
					Ipv4:                   "127.0.0.1",
					Ipv6:                   "8b34:f82b:999a:1ab5:0cad:f252:af94:bf80",
					Location:               "test",
					Name:                   "test",
					Profile:                "test",
					Slug:                   "test",
					SnapshotRunTime:        0,
					Status:                 "provisioning",
					Virtualization:         "containerd",
					WebServer:              "nginx",
					CallbackID:             "callback",
				}, nil)

				client.On("GetEvents", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
		},
		"success": {
			rd: resource.Server().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("CreateServer", ctx, mock.Anything).Once().Return(&api.Server{
					SSHPasswordAuthEnabled: true,
					WordPressLockDown:      true,
					Aliases:                []string{"test"},
					Date:                   "2022-12-22T03:54:56+03:00",
					Image:                  "test",
					Ipv4:                   "127.0.0.1",
					Ipv6:                   "8b34:f82b:999a:1ab5:0cad:f252:af94:bf80",
					Location:               "test",
					Name:                   "test",
					Profile:                "test",
					Slug:                   "test",
					SnapshotRunTime:        0,
					Status:                 "provisioning",
					Virtualization:         "containerd",
					WebServer:              "nginx",
					CallbackID:             "callback",
				}, nil)

				client.On("GetEvents", ctx, mock.Anything).Once().Return(api.Events{
					{
						Status: "finished",
					},
				}, nil)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := resource.Server().CreateContext(ctx, test.rd, config.NewCombinedConfig(&config.Config{
				ServerUpPort: 2200,
				RetryLimit:   3,
			}, client))

			assert.Equal(t, test.diags, diags)
		})
	}
}

func TestResourceWebdockServerRead(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewClientInterface(t)
	mockErr := errors.New("mock error")
	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"when get server by slug fails": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			diags: diag.Errorf("error getting server: %v", mockErr),
			mock: func() {
				client.On("GetServerBySlug", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
		},
		"success": {
			rd: resource.Server().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetServerBySlug", ctx, mock.Anything).Once().Return(&api.Server{
					SSHPasswordAuthEnabled: true,
					WordPressLockDown:      true,
					Aliases:                []string{"test"},
					Date:                   "2022-12-22T03:54:56+03:00",
					Image:                  "test",
					Ipv4:                   "83.69.106.70",
					Ipv6:                   "8b34:f82b:999a:1ab5:0cad:f252:af94:bf80",
					Location:               "test",
					Name:                   "test",
					Profile:                "test",
					Slug:                   "test",
					SnapshotRunTime:        0,
					Status:                 "provisioning",
					Virtualization:         "containerd",
					WebServer:              "nginx",
					CallbackID:             "callback",
				}, nil)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := resource.Server().ReadContext(ctx, test.rd, config.NewCombinedConfig(&config.Config{
				ServerUpPort: 2200,
				RetryLimit:   3,
			}, client))

			assert.Equal(t, test.diags, diags)
		})
	}
}
