package datasource_acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zolamk/terraform-provider-webdock/test/acctest"
)

func TestAccDataSourceServers(t *testing.T) {
	acctest.PreCheck(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:   acctest.LoadFixture(t, "testdata/servers.tf"),
				PlanOnly: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.webdock_servers.test", "servers.#"),
				),
			},
		},
	})
}
