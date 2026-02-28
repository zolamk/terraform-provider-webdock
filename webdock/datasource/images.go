package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	webdock "github.com/webdock-io/go-sdk"
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

	images, err := client.ListOSImages(webdock.ListOSImagesOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	var mappedImages []map[string]interface{}
	for _, i := range images {
		webServer := ""
		if i.WebServer != nil {
			webServer = *i.WebServer
		}
		php := ""
		if i.PHPVersion != nil {
			php = *i.PHPVersion
		}
		mappedImages = append(mappedImages, map[string]interface{}{
			"slug":        i.Slug,
			"name":        i.Name,
			"web_server":  webServer,
			"php_version": php,
		})
	}

	d.SetId("images")

	if err = d.Set("images", mappedImages); err != nil {
		return diag.Errorf("error setting images: %s", err)
	}

	return nil
}
