package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteFeatureFlagSucceedsWhenAlreadySoftDeleted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			deleted := true
			writeJSONResponse(t, w, FeatureFlag{ID: 7, Key: "pricing-test:deleted:7", Deleted: &deleted})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"type":"validation_error","code":"invalid_key","detail":"Only letters, numbers, hyphens (-) & underscores (_) are allowed.","attr":"key"}`))
	}))
	defer server.Close()

	client := newTestPosthogClient(server)
	_, err := client.DeleteFeatureFlag(context.Background(), "proj-1", "7")

	require.NoError(t, err)
}
