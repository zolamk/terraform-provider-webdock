package webdock

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWebdockImages() *schema.Resource {
	datasourceSchema := map[string]*schema.Schema{
		"images": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: false,
			Required: false,
			Elem: &schema.Resource{
				Schema: imageSchema(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: dataSourceWebdockImagesRead,
		Schema:      datasourceSchema,
	}
}

func dataSourceWebdockImagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

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
