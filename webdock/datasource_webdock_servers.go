package webdock

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/api"
)

func dataSourceWebdockServers() *schema.Resource {
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
				Schema: serverSchema(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: dataSourceWebdockServersRead,
		Schema:      datasourceSchema,
	}
}

func dataSourceWebdockServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	opts := api.GetServersParams{
		Status: d.Get("status").(string),
	}

	servers, err := client.GetServers(ctx, opts)

	if err != nil {
		return diag.Errorf("error getting servers: %s", err)
	}

	d.SetId("servers")

	if err := d.Set("servers", servers); err != nil {
		return diag.Errorf("error setting servers: %s", err)
	}

	return nil
}
