package webdock

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func imageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"slug": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image slug",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image name",
		},
		"web_server": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Web server",
		},
		"php_version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "PHP version",
		},
	}
}
