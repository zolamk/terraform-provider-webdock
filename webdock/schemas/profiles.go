package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Profile() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"slug": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Profile slug",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Profile name",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Profile description",
		},
	}
}
