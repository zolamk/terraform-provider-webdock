package webdock

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/webdock/datasource"
	"github.com/zolamk/terraform-provider-webdock/webdock/resource"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WEBDOCK_TOKEN", nil),
				Description: "The token key for API operations.",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WEBDOCK_API_URL", "https://api.webdock.io"),
				Description: "The URL to use for the Webdock API.",
			},
			"server_up_port": {
				Type:        schema.TypeInt,
				Required:    false,
				DefaultFunc: schema.EnvDefaultFunc("WEBDOCK_SERVER_UP_PORT", 22),
				Description: "The port to use when checking if the server is actually reachable.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"webdock_servers":     datasource.Servers(),
			"webdock_images":      datasource.Images(),
			"webdock_profiles":    datasource.Profiles(),
			"webdock_locations":   datasource.Locations(),
			"webdock_public_keys": datasource.PublicKeys(),
			"webdock_shell_users": datasource.ShellUsers(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"webdock_server":     resource.Server(),
			"webdock_public_key": resource.PublicKey(),
			"webdock_shell_user": resource.ShellUser(),
		},
	}

	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	config := config.Config{
		Token:            d.Get("token").(string),
		APIEndpoint:      d.Get("api_endpoint").(string),
		ServerUpPort:     d.Get("server_up_port").(int),
		TerraformVersion: terraformVersion,
	}

	return config.Client()
}
