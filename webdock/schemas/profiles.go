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
		"ram": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Profile RAM in MiB",
		},
		"disk": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Disk size in MiB",
		},
		"cpu": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "CPU model",
		},
	}
}
