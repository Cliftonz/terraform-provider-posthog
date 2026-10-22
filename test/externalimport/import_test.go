package externalimport

import (
	"testing"

	"github.com/posthog/terraform-provider/provider"
)

func TestPublicProviderPackageImportsFromExternalModule(t *testing.T) {
	if provider.New("test")() == nil {
		t.Fatal("provider.New returned nil provider")
	}
}
