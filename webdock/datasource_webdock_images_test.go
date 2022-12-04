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

func TestDataSourceWebdockImages(t *testing.T) {
	ctx := context.Background()
	client := &mocks.ClientInterface{}
	mockErr := errors.New("mock error")

	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"success": {
			rd: dataSourceWebdockImages().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetServersImages", ctx).Once().Return(api.ServerImages{
					api.ServerImage{
						Name:       "test",
						PhpVersion: "1.0",
						Slug:       "test",
						WebServer:  "test",
					},
				}, nil)
			},
		},
		"error: ": {
			rd: dataSourceWebdockImages().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetServersImages", ctx).Once().Return(nil, mockErr)
			},
			diags: diag.FromErr(errors.New("mock error")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := dataSourceWebdockImages().ReadContext(ctx, test.rd, &CombinedConfig{
				client: client,
			})

			assert.Equal(t, test.diags, diags)
		})
	}
}
