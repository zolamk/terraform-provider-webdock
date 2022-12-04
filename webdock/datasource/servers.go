package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/api"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/webdock/schemas"
)

func Servers() *schema.Resource {
	datasourceSchema := map[string]*schema.Schema{
		"status": {
			Type:        schema.TypeString,
			Default:     "all",
			Optional:    true,
			Description: "Server status (all, suspended, active)",
		},
		"servers": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: schemas.Server(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: readServers,
		Schema:      datasourceSchema,
	}
}

func readServers(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	opts := api.GetServersParams{
		Status: d.Get("status").(string),
	}

	servers, err := client.GetServers(ctx, opts)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("servers")

	if err := d.Set("servers", servers); err != nil {
		return diag.Errorf("error setting servers: %s", err)
	}

	return nil
}
