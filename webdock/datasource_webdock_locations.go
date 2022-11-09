package webdock

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWebdockLocations() *schema.Resource {
	datasourceSchema := map[string]*schema.Schema{
		"locations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: locationSchema(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: dataSourceWebdockLocationsRead,
		Schema:      datasourceSchema,
	}
}

func dataSourceWebdockLocationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

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
