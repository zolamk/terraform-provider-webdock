package webdock

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWebdockPublicKeys() *schema.Resource {
	datasourceSchema := map[string]*schema.Schema{
		"public_keys": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: publicKeySchema(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: dataSourceWebdockPublicKeysRead,
		Schema:      datasourceSchema,
	}
}

func dataSourceWebdockPublicKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	publicKeys, err := client.GetPublicKeys(ctx)

	if err != nil {
		return diag.Errorf("error getting public keys: %s", err)
	}

	d.SetId("public_keys")

	if err = d.Set("public_keys", publicKeys); err != nil {
		return diag.Errorf("error setting public keys: %s", err)
	}

	return nil
}
