package datasource_acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zolamk/terraform-provider-webdock/test/acctest"
)

// TestAccDataSourceProfiles requires a known location to be specified.
// Use WEBDOCK_LOCATION_ID env var (e.g. "fi") or update the config inline.
func TestAccDataSourceProfiles(t *testing.T) {
	acctest.PreCheck(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:   acctest.LoadFixture(t, "testdata/profiles.tf"),
				PlanOnly: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.webdock_profiles.test", "profiles.#"),
					resource.TestCheckResourceAttrSet("data.webdock_profiles.test", "profiles.0.slug"),
					resource.TestCheckResourceAttrSet("data.webdock_profiles.test", "profiles.0.name"),
				),
			},
		},
	})
}
