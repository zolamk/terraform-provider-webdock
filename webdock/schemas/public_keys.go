package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func PublicKey() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "PublicKey ID",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "PublicKey name",
		},
		"key": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "PublicKey content",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "PublicKey creation datetime",
		},
	}
}
