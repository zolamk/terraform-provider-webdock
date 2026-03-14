package resource_acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zolamk/terraform-provider-webdock/test/acctest"
)

// TestAccResourceServer_plan verifies the webdock_server resource produces a
// valid plan. Nothing is created or destroyed in the account.
func TestAccResourceServer_plan(t *testing.T) {
	acctest.PreCheck(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config:   acctest.LoadFixture(t, "testdata/server.tf"),
				// Verify that all required computed attributes appear in the plan.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("webdock_server.test", "name", "acc-test-server"),
					resource.TestCheckResourceAttr("webdock_server.test", "location_id", "fi"),
					resource.TestCheckResourceAttr("webdock_server.test", "profile_slug", "webdock-standard-2vcpu-2gb"),
					resource.TestCheckResourceAttr("webdock_server.test", "image_slug", "ubuntu-2404-noble"),
				),
			},
		},
	})
}

// TestAccResourceServer_withPublicKey_plan tests a full server + public-key
// + shell-user composition at plan time to ensure provider-level dependencies
// are expressed correctly.
func TestAccResourceServer_withPublicKey_plan(t *testing.T) {
	acctest.PreCheck(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config:   acctest.LoadFixture(t, "testdata/server_with_public_key.tf"),
			},
		},
	})
}
