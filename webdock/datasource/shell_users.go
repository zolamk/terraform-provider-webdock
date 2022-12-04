package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/webdock/schemas"
)

func ShellUsers() *schema.Resource {
	datasourceSchema := map[string]*schema.Schema{
		"server_slug": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"shell_users": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: schemas.ShellUser(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: readShellUsers,
		Schema:      datasourceSchema,
	}
}

func readShellUsers(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	shellUsers, err := client.GetShellUsers(ctx, d.Get("server_slug").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("shell_users")

	if err = d.Set("shell_users", shellUsers); err != nil {
		return diag.Errorf("error setting shell users: %s", err)
	}

	return nil
}
