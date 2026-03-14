package resource_acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zolamk/terraform-provider-webdock/test/acctest"
)

// TestAccResourcePublicKey_plan verifies the webdock_public_key resource
// produces a valid plan without applying any changes to the account.
func TestAccResourcePublicKey_plan(t *testing.T) {
	acctest.PreCheck(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Only run plan – do not apply, so no key is added to the account.
				PlanOnly: true,
				Config:   acctest.LoadFixture(t, "testdata/public_key.tf"),
			},
		},
	})
}
