package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	webdock "github.com/webdock-io/go-sdk"
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

	opts := webdock.ListServerOptions{
		Status: webdock.ListServersQuery(d.Get("status").(string)),
	}

	servers, err := client.ListServer(opts)

	if err != nil {
		return diag.FromErr(err)
	}

	var mappedServers []map[string]interface{}
	for _, s := range servers {
		mappedServers = append(mappedServers, map[string]interface{}{
			"aliases":                   s.Aliases,
			"created_at":                s.Date,
			"image_slug":                s.Image,
			"ipv4":                      s.IPv4,
			"ipv6":                      s.IPv6,
			"location_id":               s.Location,
			"name":                      s.Name,
			"profile_slug":              s.Profile,
			"slug":                      s.Slug,
			"snapshot_runtime":          s.SnapshotRunTime,
			"status":                    string(s.Status),
			"wordpress_lockdown":        s.WordPressLockDown,
			"webserver":                 string(s.WebServer),
			"ssh_password_auth_enabled": s.SSHPasswordAuthEnabled,
			"virtualization":            string(s.Virtualization),
		})
	}

	d.SetId("servers")

	if err := d.Set("servers", mappedServers); err != nil {
		return diag.Errorf("error setting servers: %s", err)
	}

	return nil
}
