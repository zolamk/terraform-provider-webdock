// Package acctest provides shared helpers for acceptance (integration) tests.
//
// Tests in this package require a real Webdock API token set in the
// WEBDOCK_TOKEN environment variable. If the variable is absent the test is
// skipped automatically, so the suite is always safe to run in environments
// without credentials.
package acctest

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zolamk/terraform-provider-webdock/webdock"
)

// ProviderFactories is the map of provider factories required by
// resource.TestCase when using the SDKv2 framework.
var ProviderFactories = map[string]func() (*schema.Provider, error){
	"webdock": func() (*schema.Provider, error) {
		return webdock.Provider(), nil
	},
}

// PreCheck skips the test when WEBDOCK_TOKEN is not set, and additionally
// validates that any extra required variables are present.
func PreCheck(t *testing.T, extra ...string) {
	t.Helper()

	if os.Getenv("WEBDOCK_TOKEN") == "" {
		t.Skip("WEBDOCK_TOKEN not set – skipping acceptance test")
	}

	for _, env := range extra {
		if os.Getenv(env) == "" {
			t.Skipf("%s not set – skipping acceptance test", env)
		}
	}
}

// ProviderConfig returns a minimal provider block that can be embedded in
// any test configuration. The token is read from the environment through the
// provider's own DefaultFunc, so no value needs to be hard-coded here.
const ProviderConfig = `
provider "webdock" {}
`

// CheckFunc is a convenience alias for resource.TestCheckFunc.
type CheckFunc = resource.TestCheckFunc

// ComposeChecks is a convenience alias for resource.ComposeAggregateTestCheckFunc.
var ComposeChecks = resource.ComposeAggregateTestCheckFunc

// LoadFixture reads a file from the testdata directory and prepends the ProviderConfig.
func LoadFixture(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %s", path, err)
	}
	return ProviderConfig + string(b)
}
