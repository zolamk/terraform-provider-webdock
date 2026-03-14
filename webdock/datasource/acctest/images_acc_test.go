package datasource_acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zolamk/terraform-provider-webdock/test/acctest"
)

func TestAccDataSourceImages(t *testing.T) {
	acctest.PreCheck(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:   acctest.LoadFixture(t, "testdata/images.tf"),
				PlanOnly: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.webdock_images.test", "images.#"),
					resource.TestCheckResourceAttrSet("data.webdock_images.test", "images.0.slug"),
					resource.TestCheckResourceAttrSet("data.webdock_images.test", "images.0.name"),
				),
			},
		},
	})
}
