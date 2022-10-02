package webdock

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func locationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location ID",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location name",
		},
		"city": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location city",
		},
		"country": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location string",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location description",
		},
		"icon": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location icon",
		},
	}
}
