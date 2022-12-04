package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/webdock/schemas"
)

func Images() *schema.Resource {
	datasourceSchema := map[string]*schema.Schema{
		"images": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: false,
			Required: false,
			Elem: &schema.Resource{
				Schema: schemas.Image(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: readImages,
		Schema:      datasourceSchema,
	}
}

func readImages(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	images, err := client.GetServersImages(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("images")

	if err = d.Set("images", images); err != nil {
		return diag.Errorf("error setting images: %s", err)
	}

	return nil
}
