package datasource_test

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
	"github.com/zolamk/terraform-provider-webdock/webdock/datasource"
)

func TestDataSourceWebdockShellUsersRead(t *testing.T) {
	ctx := context.Background()
	client := &mocks.ClientInterface{}
	mockErr := errors.New("mock error")

	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"success": {
			rd: datasource.ShellUsers().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetShellUsers", ctx, mock.Anything).Once().Return(api.ShellUsers{
					api.ShellUser{
						ID:         json.Number("1"),
						Username:   "test",
						Password:   "test",
						Group:      "test",
						Shell:      "test",
						PublicKeys: api.PublicKeys{},
						Created:    "04/01/2022 06:36:01",
					},
				}, nil)
			},
		},
		"error: ": {
			rd: datasource.ShellUsers().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetShellUsers", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
			diags: diag.FromErr(errors.New("mock error")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := datasource.ShellUsers().ReadContext(ctx, test.rd, config.NewCombinedConfig(&config.Config{
				ServerUpPort: 2200,
			}, client))

			assert.Equal(t, test.diags, diags)
		})
	}
}
