package webdock

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		},
		DataSourcesMap: map[string]*schema.Resource{
			"webdock_account":  datasourceWebdockAccount(),
			"webdock_servers":  dataSourceWebdockServers(),
			"webdock_images":   dataSourceWebdockImages(),
			"webdock_profiles": dataSourceWebdockProfiles(),
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
	config := Config{
		Token:            d.Get("token").(string),
		APIEndpoint:      d.Get("api_endpoint").(string),
		TerraformVersion: terraformVersion,
	}

	return config.Client()
}
