package datasource_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	webdock "github.com/webdock-io/go-sdk"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/test/mocks"
	"github.com/zolamk/terraform-provider-webdock/webdock/datasource"
)

func TestDataSourceWebdockServers(t *testing.T) {
	ctx := context.Background()
	client := &mocks.ClientInterface{}
	mockErr := errors.New("mock error")

	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"success": {
			rd: datasource.Servers().Data(&terraform.InstanceState{
				Attributes: map[string]string{
					"status": "all",
				},
			}),
			mock: func() {
				opts := webdock.ListServerOptions{
					Status: webdock.ListServersQuery("all"),
				}
				client.On("ListServer", opts).Once().Return(webdock.ListServers{
					{
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
			},
		},
		"error:": {
			rd: datasource.Servers().Data(&terraform.InstanceState{
				Attributes: map[string]string{
					"status": "all",
				},
			}),
			mock: func() {
				opts := webdock.ListServerOptions{
					Status: webdock.ListServersQuery("all"),
				}
				client.On("ListServer", opts).Once().Return(nil, mockErr)
			},
			diags: diag.FromErr(errors.New("mock error")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := datasource.Servers().ReadContext(ctx, test.rd, config.NewCombinedConfig(&config.Config{
				ServerUpPort: 2200,
			}, client))

			assert.Equal(t, test.diags, diags)
		})
	}
}
