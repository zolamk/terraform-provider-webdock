package webdock

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceWebdockShellUsers() *schema.Resource {
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
				Schema: shellUserSchema(),
			},
		},
	}

	return &schema.Resource{
		ReadContext: dataSourceWebdockShellUsersRead,
		Schema:      datasourceSchema,
	}
}

func dataSourceWebdockShellUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	shellUsers, err := client.GetShellUsers(ctx, d.Get("server_slug").(string))
	if err != nil {
		return diag.Errorf("error getting shell users: %s", err)
	}

	d.SetId("shell_users")

	if err = d.Set("shell_users", shellUsers); err != nil {
		return diag.Errorf("error setting shell users: %s", err)
	}

	return nil
}
