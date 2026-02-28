package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	webdock "github.com/webdock-io/go-sdk"
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

	opts := webdock.ListPossibleProfilesInLocationOptions{
		LocationID: d.Get("location_id").(string),
	}

	profiles, err := client.ListPossibleProfilesInLocation(opts)
	if err != nil {
		return diag.FromErr(err)
	}

	var mappedProfiles []map[string]interface{}
	for _, p := range profiles {
		mappedProfiles = append(mappedProfiles, map[string]interface{}{
			"slug":        p.Slug,
			"name":        p.Name,
			"description": p.Description,
		})
	}

	d.SetId("profiles")

	if err = d.Set("profiles", mappedProfiles); err != nil {
		return diag.Errorf("error setting profiles: %s", err)
	}

	return nil
}
