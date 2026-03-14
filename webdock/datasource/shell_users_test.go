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

func TestDataSourceWebdockShellUsers(t *testing.T) {
	ctx := context.Background()
	client := &mocks.ClientInterface{}
	mockErr := errors.New("mock error")

	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"success": {
			rd: datasource.ShellUsers().Data(&terraform.InstanceState{
				Attributes: map[string]string{
					"server_slug": "test",
				},
			}),
			mock: func() {
				opts := webdock.ListServerShellUserOptions{
					ServerSlug: "test",
				}
				client.On("ListServerShellUser", opts).Once().Return([]webdock.ShellUser{
					{
						ID:         1,
						Username:   "test",
						Group:      "test",
						Shell:      "test",
						PublicKeys: []webdock.PublicKeyDTO{{ID: 1}},
						Created:    "2022-12-22T03:54:56+03:00",
					},
				}, nil)
			},
		},
		"error:": {
			rd: datasource.ShellUsers().Data(&terraform.InstanceState{
				Attributes: map[string]string{
					"server_slug": "test",
				},
			}),
			mock: func() {
				opts := webdock.ListServerShellUserOptions{
					ServerSlug: "test",
				}
				client.On("ListServerShellUser", opts).Once().Return(nil, mockErr)
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
