package webdock

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func serverSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aliases": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "Server description (what's installed here?) as entered by admin in Server Metadata",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Creation date/time",
		},
		"image_slug": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Server image",
		},
		"ipv4": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv4 address",
		},
		"ipv6": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IPv6 address",
		},
		"location_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Location ID of the server",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Server name",
		},
		"next_action_date": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Next action date/time as entered by admin in Server Metadata",
		},
		"notes": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Notes as entered by admin in Server Metadata",
		},
		"profile_slug": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Server profile",
		},
		"slug": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Server slug",
		},
		"snapshot_runtime": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Last knows snapshot runtime (seconds)",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Server status",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Server description (what's installed here?) as entered by admin in Server Metadata",
		},
		"wordpress_lockdown": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Whether WordPress is in lockdown mode",
		},
		"webserver": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Webserver type (apache, nginx, none)",
		},
		"ssh_password_auth_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Whether SSH password authentication is enabled",
		},
		"virtualization": {
			Type:        schema.TypeString,
			Default:     "container",
			Optional:    true,
			Description: "Virtualization type for your new server. container means the server will be a Webdock LXD VPS and kvm means it will be a KVM Virtual machine. If you specify a snapshotId in the request, the server type from which the snapshot belongs much match the virtualization selected. Reason being that KVM images are incompatible with LXD images and vice-versa.",
		},
	}
}
