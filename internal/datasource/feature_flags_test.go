package datasource

import (
	"context"
	"testing"

	"github.com/posthog/terraform-provider/internal/httpclient"
	"github.com/posthog/terraform-provider/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeatureFlagsToList_DropsDeletedFlagsAndSortsByKey(t *testing.T) {
	flags := []httpclient.FeatureFlag{
		{ID: 2, Key: "zeta", Name: util.StringPtr("Zeta"), Active: util.BoolPtr(true)},
		{ID: 3, Key: "gone", Deleted: util.BoolPtr(true)},
		{ID: 1, Key: "alpha", Name: util.StringPtr("  "), Active: util.BoolPtr(false), Archived: util.BoolPtr(true)},
	}

	list, diags := featureFlagsToList(context.Background(), flags)
	require.False(t, diags.HasError())

	var got []featureFlagSummaryModel
	require.False(t, list.ElementsAs(context.Background(), &got, false).HasError())
	require.Len(t, got, 2)

	assert.Equal(t, "alpha", got[0].Key.ValueString())
	assert.Equal(t, int64(1), got[0].ID.ValueInt64())
	assert.True(t, got[0].Name.IsNull(), "whitespace-only name should map to null")
	assert.False(t, got[0].Active.ValueBool())
	assert.True(t, got[0].Archived.ValueBool())

	assert.Equal(t, "zeta", got[1].Key.ValueString())
	assert.Equal(t, "Zeta", got[1].Name.ValueString())
	assert.True(t, got[1].Active.ValueBool())
	assert.False(t, got[1].Archived.ValueBool(), "a missing archived field means not archived")
}

func TestFeatureFlagsToList_NoFlagsIsAnEmptyListNotNull(t *testing.T) {
	list, diags := featureFlagsToList(context.Background(), nil)
	require.False(t, diags.HasError())

	assert.False(t, list.IsNull(), "an empty project should give an empty list, so for_each over it works")
	assert.Empty(t, list.Elements())
}

func TestMatchProjectByName(t *testing.T) {
	projects := []httpclient.Project{
		{ID: 1, Name: util.StringPtr("Production")},
		{ID: 2, Name: util.StringPtr(" staging ")},
		{ID: 3, Name: nil},
		{ID: 4, Name: util.StringPtr("dup")},
		{ID: 5, Name: util.StringPtr("DUP")},
	}

	id, found, err := matchProjectByName(projects, "production")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "1", id)

	id, found, err = matchProjectByName(projects, "staging")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "2", id)

	_, found, err = matchProjectByName(projects, "preview-123")
	require.NoError(t, err)
	assert.False(t, found, "an unmatched name is not an error")

	_, _, err = matchProjectByName(projects, "dup")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "IDs 4, 5")
}
