package results

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleVisualizer_LoadAndSelect(t *testing.T) {
	v := NewSimpleVisualizer()
	rs := ResultSet{
		Columns: []string{"id", "name"},
		Rows:    []Row{{"id": 1, "name": "a"}, {"id": 2, "name": "b"}},
	}

	require.NoError(t, v.Load(rs))
	if r := v.SelectedRow(); r != nil {
		t.Fatalf("expected no selection, got %v", r)
	}

	v.Select(1)
	r := v.SelectedRow()
	require.NotNil(t, r)
	require.Equal(t, 2, r["id"])
}
