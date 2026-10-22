package provider_test

import (
	"context"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/posthog/terraform-provider/provider"
)

func TestNewIsImportableFromOutsideModule(t *testing.T) {
	var resp fwprovider.MetadataResponse
	provider.New("test")().Metadata(context.Background(), fwprovider.MetadataRequest{}, &resp)

	if resp.TypeName != "posthog" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "posthog")
	}
}
