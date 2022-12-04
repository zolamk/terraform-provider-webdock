package webdock

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
	"github.com/zolamk/terraform-provider-webdock/test/mocks"
)

func TestDataSourceWebdockPublicKeysRead(t *testing.T) {
	ctx := context.Background()
	client := &mocks.ClientInterface{}
	mockErr := errors.New("mock error")

	tests := map[string]struct {
		rd    *schema.ResourceData
		diags diag.Diagnostics
		mock  func()
	}{
		"success": {
			rd: dataSourceWebdockPublicKeys().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetPublicKeys", ctx, mock.Anything).Once().Return(api.PublicKeys{
					api.PublicKey{
						Id:      json.Number("1"),
						Created: "02/03/2022 20:37:27",
						Name:    "test",
						Key:     "public key content",
					},
				}, nil)
			},
		},
		"error: ": {
			rd: dataSourceWebdockPublicKeys().Data(&terraform.InstanceState{}),
			mock: func() {
				client.On("GetPublicKeys", ctx, mock.Anything).Once().Return(nil, mockErr)
			},
			diags: diag.FromErr(errors.New("mock error")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.mock()

			diags := dataSourceWebdockPublicKeys().ReadContext(ctx, test.rd, &CombinedConfig{
				client: client,
			})

			assert.Equal(t, test.diags, diags)
		})
	}
}
