package dbconn_test

import (
	"testing"

	"github.com/panda-coder/GoScratch/pkg/dbconn"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	m := dbconn.NewManager()
	require.NotNil(t, m)

	err := m.Add("local", "sqlite://:memory:")
	require.NoError(t, err)

	ci, ok := m.Get("local")
	require.True(t, ok)
	require.Equal(t, "local", ci.Name)
	require.Equal(t, "sqlite://:memory:", ci.DSN)

	list := m.List()
	require.Len(t, list, 1)

	m.Remove("local")
	_, ok = m.Get("local")
	require.False(t, ok)
}
