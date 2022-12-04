package resource_test

import (
	"context"
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

func TestResourceWebdockServerCreate(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewClientInterface(t)
	mockErr := errors.New("mock error")
	tests := map[string]struct {
		rd    *schema.ResourceData
		meta  interface{}
		diags diag.Diagnostics
		mock  func()
	}{
		"when create server fails": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			meta:  config.NewCombinedConfig(nil, client),
			diags: diag.FromErr(mockErr),
			mock: func() {
				client.On("CreateServer", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
		},
		"when wait for action fails": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			meta:  config.NewCombinedConfig(nil, client),
			diags: diag.Errorf("server (test) create event (callback) errored: %v", mockErr),
			mock: func() {
				client.On("CreateServer", ctx, mock.Anything).Once().Return(&api.Server{
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

				client.On("GetEvents", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
		},
		"when create event can't be found": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			meta:  config.NewCombinedConfig(nil, client),
			diags: diag.Errorf("unable to find server (test) create event"),
			mock: func() {
				client.On("CreateServer", ctx, mock.Anything).Once().Return(&api.Server{
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
				}, nil)
			},
		},
		"success": {
			rd:   resource.Server().Data(&terraform.InstanceState{}),
			meta: config.NewCombinedConfig(nil, client),
			mock: func() {
				client.On("CreateServer", ctx, mock.Anything).Once().Return(&api.Server{
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

			diags := resource.Server().CreateContext(ctx, test.rd, test.meta)

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
		meta  interface{}
		diags diag.Diagnostics
		mock  func()
	}{
		"when get server by slug fails": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			meta:  config.NewCombinedConfig(nil, client),
			diags: diag.Errorf("error getting server: %v", mockErr),
			mock: func() {
				client.On("GetServerBySlug", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
		},
		"success": {
			rd:   resource.Server().Data(&terraform.InstanceState{}),
			meta: config.NewCombinedConfig(nil, client),
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

			diags := resource.Server().ReadContext(ctx, test.rd, test.meta)

			assert.Equal(t, test.diags, diags)
		})
	}
}
