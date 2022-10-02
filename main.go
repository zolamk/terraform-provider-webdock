package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/zolamk/terraform-provider-webdock/webdock"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: webdock.Provider,
	})
}
