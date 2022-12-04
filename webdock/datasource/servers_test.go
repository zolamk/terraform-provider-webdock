package datasource_test

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
	"github.com/zolamk/terraform-provider-webdock/webdock/datasource"
)

func TestDataSourceWebdockServersRead(t *testing.T) {
	ctx := context.Background()
	client := &mocks.ClientInterface{}
	mockErr := errors.New("mock error")

	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"success": {
			rd: datasource.Servers().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetServers", ctx, mock.Anything).Once().Return(api.Servers{
					api.Server{
						SSHPasswordAuthEnabled: true,
						WordPressLockDown:      false,
						Aliases:                []string{"alias1", "alias2"},
						Date:                   "19/10/2022 03:12:13",
						Image:                  "test",
						Ipv4:                   "149.57.225.5",
						Ipv6:                   "b946:997f:cc88:0251:e7dd:dd2e:c3ef:3764",
						Location:               "test",
						Name:                   "test",
						Profile:                "test",
						Slug:                   "test",
						SnapshotRunTime:        10,
						Status:                 "test",
						Virtualization:         "test",
						WebServer:              "test",
					},
				}, nil)
			},
		},
		"error: ": {
			rd: datasource.Servers().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetServers", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
			diags: diag.FromErr(errors.New("mock error")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := datasource.Servers().ReadContext(ctx, test.rd, config.NewCombinedConfig(nil, client))

			assert.Equal(t, test.diags, diags)
		})
	}
}
