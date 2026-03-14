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
	webdock "github.com/webdock-io/go-sdk"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/test/mocks"
	"github.com/zolamk/terraform-provider-webdock/webdock/resource"
)

func TestResourceWebdockServerCreate(t *testing.T) {
	ctx := context.Background()
	client := mocks.NewClientInterface(t)
	mockErr := errors.New("mock error")
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.Nil(t, err)
	defer l.Close()
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	testPort := l.Addr().(*net.TCPAddr).Port

	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"when create server fails": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			diags: diag.FromErr(mockErr),
			mock: func() {
				client.On("CreateServerFromImage", mock.Anything).Once().Return(webdock.CreatedServer{}, mockErr)
			},
		},
		"when create server fails with too many server error": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			diags: nil,
			mock: func() {
				// The provider retries once on "You are creating too many servers". So we mock it twice.
				client.On("CreateServerFromImage", mock.Anything).Twice().Return(webdock.CreatedServer{}, errors.New("You are creating too many servers in too short of a timespan. Please wait a while and try again a bit later."))

				client.On("CreateServerFromImage", mock.Anything).Once().Return(webdock.CreatedServer{
					CallbackID: "callback",
					Server: webdock.Server{
						SSHPasswordAuthEnabled: true,
						WordPressLockDown:      true,
						Aliases:                []string{"test"},
						Date:                   "2022-12-22T03:54:56+03:00",
						Image:                  "test",
						IPv4:                   "127.0.0.1",
						IPv6:                   "8b34:f82b:999a:1ab5:0cad:f252:af94:bf80",
						Location:               "test",
						Name:                   "test",
						Profile:                "test",
						Slug:                   "test",
						SnapshotRunTime:        0,
						Status:                 "provisioning",
						Virtualization:         "containerd",
						WebServer:              "nginx",
					},
				}, nil)

				client.On("ListEvents", mock.Anything).Once().Return(webdock.ListEventsResponse{
					Events: []webdock.EventDTO{
						{Status: "finished"},
					},
				}, nil)
			},
		},
		"when wait for action fails": {
			rd:    resource.Server().Data(&terraform.InstanceState{}),
			diags: diag.Errorf("server (test) create event (callback) errored: %v", mockErr),
			mock: func() {
				client.On("CreateServerFromImage", mock.Anything).Once().Return(webdock.CreatedServer{
					CallbackID: "callback",
					Server: webdock.Server{
						SSHPasswordAuthEnabled: true,
						WordPressLockDown:      true,
						Aliases:                []string{"test"},
						Date:                   "2022-12-22T03:54:56+03:00",
						Image:                  "test",
						IPv4:                   "127.0.0.1",
						IPv6:                   "8b34:f82b:999a:1ab5:0cad:f252:af94:bf80",
						Location:               "test",
						Name:                   "test",
						Profile:                "test",
						Slug:                   "test",
						SnapshotRunTime:        0,
						Status:                 "provisioning",
						Virtualization:         "containerd",
						WebServer:              "nginx",
					},
				}, nil)

				client.On("ListEvents", mock.Anything).Once().Return(webdock.ListEventsResponse{}, mockErr)
			},
		},
		"success": {
			rd: resource.Server().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("CreateServerFromImage", mock.Anything).Once().Return(webdock.CreatedServer{
					CallbackID: "callback",
					Server: webdock.Server{
						SSHPasswordAuthEnabled: true,
						WordPressLockDown:      true,
						Aliases:                []string{"test"},
						Date:                   "2022-12-22T03:54:56+03:00",
						Image:                  "test",
						IPv4:                   "127.0.0.1",
						IPv6:                   "8b34:f82b:999a:1ab5:0cad:f252:af94:bf80",
						Location:               "test",
						Name:                   "test",
						Profile:                "test",
						Slug:                   "test",
						SnapshotRunTime:        0,
						Status:                 "provisioning",
						Virtualization:         "containerd",
						WebServer:              "nginx",
					},
				}, nil)

				client.On("ListEvents", mock.Anything).Once().Return(webdock.ListEventsResponse{
					Events: []webdock.EventDTO{
						{Status: "finished"},
					},
				}, nil)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := resource.Server().CreateContext(ctx, test.rd, config.NewCombinedConfig(&config.Config{
				ServerUpPort: testPort,
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
				client.On("GetServerBySlug", mock.Anything).Once().Return(webdock.Server{}, mockErr)
			},
		},
		"success": {
			rd: resource.Server().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetServerBySlug", mock.Anything).Once().Return(webdock.Server{
					SSHPasswordAuthEnabled: true,
					WordPressLockDown:      true,
					Aliases:                []string{"test"},
					Date:                   "2022-12-22T03:54:56+03:00",
					Image:                  "test",
					IPv4:                   "83.69.106.70",
					IPv6:                   "8b34:f82b:999a:1ab5:0cad:f252:af94:bf80",
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
