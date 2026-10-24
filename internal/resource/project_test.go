package resource

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/posthog/terraform-provider/internal/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectDeleteSucceedsWhenAlreadyBeingDeleted(t *testing.T) {
	const orgID = "01890a5c-7a3b-0000-8f1e-5c2d9e1b4a77"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/organizations/"+orgID+"/projects/42/", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"type":"validation_error","code":"invalid","detail":"This project is already being deleted.","attr":null}`))
	}))
	defer server.Close()

	client := httpclient.NewClient(server.Client(), server.URL, "test-key", "test")
	model := ProjectTFModel{}
	model.ID = types.Int64Value(42)
	model.OrganizationID = types.StringValue(orgID)

	_, err := ProjectOps{}.Delete(context.Background(), client, model)

	require.NoError(t, err)
}
