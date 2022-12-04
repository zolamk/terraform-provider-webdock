package webdock

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/zolamk/terraform-provider-webdock/api"
	"github.com/zolamk/terraform-provider-webdock/test/mocks"
)

func TestDataSourceWebdockLocationsRead(t *testing.T) {
	ctx := context.Background()
	client := &mocks.ClientInterface{}
	mockErr := errors.New("mock error")

	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"success": {
			rd: dataSourceWebdockLocations().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetServersLocations", ctx).Once().Return(api.ServerLocations{
					api.ServerLocation{
						City:        "test",
						Country:     "test",
						Description: "test",
						Icon:        "test",
						ID:          "test",
						Name:        "test",
					},
				}, nil)
			},
		},
		"error: ": {
			rd: dataSourceWebdockLocations().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetServersLocations", ctx).Once().Return(nil, mockErr)
			},
			diags: diag.FromErr(errors.New("mock error")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := dataSourceWebdockLocations().ReadContext(ctx, test.rd, &CombinedConfig{
				client: client,
			})

			assert.Equal(t, test.diags, diags)
		})
	}
}
