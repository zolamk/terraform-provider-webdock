package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zolamk/terraform-provider-webdock/api"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/webdock/schemas"
)

func Profiles() *schema.Resource {
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
				Schema: schemas.Profile(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: readProfiles,
		Schema:      datasourceSchema,
	}
}

func readProfiles(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	opts := api.GetServersProfilesParams{
		LocationId: d.Get("location_id").(string),
	}

	profiles, err := client.GetServersProfiles(ctx, opts)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("profiles")

	if err = d.Set("profiles", profiles); err != nil {
		return diag.Errorf("error setting profiles: %s", err)
	}

	return nil
}
