package db

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewStore(t *testing.T) {
	store := NewStore(testDB)

	require.NotEmpty(t, store)
}
