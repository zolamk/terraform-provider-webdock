package datasource_acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zolamk/terraform-provider-webdock/test/acctest"
)

func TestAccDataSourcePublicKeys(t *testing.T) {
	acctest.PreCheck(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:   acctest.LoadFixture(t, "testdata/public_keys.tf"),
				PlanOnly: true,
				// The account may have zero keys – we only assert the attribute
				// itself is present (even if empty).
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.webdock_public_keys.test", "public_keys.#"),
				),
			},
		},
	})
}
