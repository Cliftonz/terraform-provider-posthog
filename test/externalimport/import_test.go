package externalimport

import (
	"testing"

	"github.com/posthog/terraform-provider/provider"
)

func TestPublicProviderPackageImportsFromExternalModule(t *testing.T) {
	_ = provider.New("test")
}
