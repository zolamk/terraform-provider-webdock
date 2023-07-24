package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ShellUser() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "shell user id",
		},
		"server_slug": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "shell user server slug",
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "shell user username",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Sensitive:   true,
			Description: "shell user password",
		},
		"group": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "shell user group",
			Default:     "sudo",
		},
		"shell": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "shell user shell",
			Default:     "/bin/bash",
		},
		"public_keys": {
			Type:        schema.TypeList,
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Description: "shell user public keys",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "shell user creation datetime",
		},
	}
}
