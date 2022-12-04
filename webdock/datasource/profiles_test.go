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

func TestDataSourceWebdockProfilesRead(t *testing.T) {
	ctx := context.Background()
	client := &mocks.ClientInterface{}
	mockErr := errors.New("mock error")

	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"success": {
			rd: datasource.Profiles().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetServersProfiles", ctx, mock.Anything).Once().Return(api.ServerProfiles{
					api.ServerProfile{
						CPU: api.CPU{
							Cores:   4,
							Threads: 8,
						},
						Disk: 10,
						Name: "test",
						Price: api.Price{
							Amount:   10,
							Currency: "USD",
						},
						RAM:  10,
						Slug: "test",
					},
				}, nil)
			},
		},
		"error: ": {
			rd: datasource.Profiles().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetServersProfiles", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
			diags: diag.FromErr(errors.New("mock error")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := datasource.Profiles().ReadContext(ctx, test.rd, config.NewCombinedConfig(nil, client))

			assert.Equal(t, test.diags, diags)
		})
	}
}
