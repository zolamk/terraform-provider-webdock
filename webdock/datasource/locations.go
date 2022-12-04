package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/webdock/schemas"
)

func Locations() *schema.Resource {
	datasourceSchema := map[string]*schema.Schema{
		"locations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: schemas.Location(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: readLocations,
		Schema:      datasourceSchema,
	}
}

func readLocations(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	locations, err := client.GetServersLocations(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("locations")

	if err = d.Set("locations", locations); err != nil {
		return diag.Errorf("error setting locations: %s", err)
	}

	return nil
}
