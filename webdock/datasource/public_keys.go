package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	webdock "github.com/webdock-io/go-sdk"
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

	publicKeys, err := client.ListAccountPublicKeys(webdock.ListAccountPublicKeysOptions{})

	if err != nil {
		return diag.FromErr(err)
	}

	var mappedPublicKeys []map[string]interface{}
	for _, p := range publicKeys {
		mappedPublicKeys = append(mappedPublicKeys, map[string]interface{}{
			"id":         fmt.Sprintf("%d", p.ID),
			"name":       p.Name,
			"key":        p.Key,
			"created_at": p.Created,
		})
	}

	d.SetId("public_keys")

	if err = d.Set("public_keys", mappedPublicKeys); err != nil {
		return diag.Errorf("error setting public keys: %s", err)
	}

	return nil
}
