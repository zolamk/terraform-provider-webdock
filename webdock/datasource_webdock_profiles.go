package webdock

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zolamk/terraform-provider-webdock/api"
)

func dataSourceWebdockProfiles() *schema.Resource {
	datasourceSchema := map[string]*schema.Schema{
		"location_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"profiles": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: profileSchema(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: dataSourceWebdockProfilesRead,
		Schema:      datasourceSchema,
	}
}

func dataSourceWebdockProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	opts := &api.GetServersProfilesParams{
		LocationId: d.Get("location_id").(string),
	}

	profiles, err := client.GetServersProfiles(ctx, opts)

	if err != nil {
		return diag.Errorf("error getting profiles: %s", err)
	}

	d.SetId("profiles")

	if err = d.Set("profiles", profiles); err != nil {
		return diag.Errorf("error setting profiles: %s", err)
	}

	return nil
}
