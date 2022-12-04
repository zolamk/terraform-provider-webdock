package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/webdock/schemas"
)

func PublicKeys() *schema.Resource {
	datasourceSchema := map[string]*schema.Schema{
		"public_keys": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: schemas.PublicKey(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: readPublicKeys,
		Schema:      datasourceSchema,
	}
}

func readPublicKeys(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	publicKeys, err := client.GetPublicKeys(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("public_keys")

	if err = d.Set("public_keys", publicKeys); err != nil {
		return diag.Errorf("error setting public keys: %s", err)
	}

	return nil
}
