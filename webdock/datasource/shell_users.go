package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	webdock "github.com/webdock-io/go-sdk"
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

	opts := webdock.ListServerShellUserOptions{
		ServerSlug: d.Get("server_slug").(string),
	}

	shellUsers, err := client.ListServerShellUser(opts)
	if err != nil {
		return diag.FromErr(err)
	}

	var mappedShellUsers []map[string]interface{}
	for _, s := range shellUsers {
		var publicKeys []interface{}
		for _, key := range s.PublicKeys {
			publicKeys = append(publicKeys, int64(key.ID))
		}

		mappedShellUsers = append(mappedShellUsers, map[string]interface{}{
			"id":          fmt.Sprintf("%d", s.ID),
			"server_slug": d.Get("server_slug").(string),
			"username":    s.Username,
			"group":       s.Group,
			"shell":       s.Shell,
			"public_keys": publicKeys,
			"created_at":  s.Created,
		})
	}

	d.SetId("shell_users")

	if err = d.Set("shell_users", mappedShellUsers); err != nil {
		return diag.Errorf("error setting shell users: %s", err)
	}

	return nil
}
