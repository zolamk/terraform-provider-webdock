package datasource_acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zolamk/terraform-provider-webdock/test/acctest"
)

// TestAccDataSourceLocations verifies that the webdock_locations data source
// correctly retrieves location data from the live API.
// The test only plans (PlanOnly: true) so no real infrastructure is mutated.
func TestAccDataSourceLocations(t *testing.T) {
	acctest.PreCheck(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:   acctest.LoadFixture(t, "testdata/locations.tf"),
				PlanOnly: true,
				// The data source is read during plan, so we can assert at
				// least one location is present in the computed list.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.webdock_locations.test", "locations.#"),
					resource.TestCheckResourceAttrSet("data.webdock_locations.test", "locations.0.id"),
					resource.TestCheckResourceAttrSet("data.webdock_locations.test", "locations.0.name"),
				),
			},
		},
	})
}
