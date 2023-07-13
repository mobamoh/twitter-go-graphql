package test_helper

import (
	"context"
	"github.com/mobamoh/twitter-go-graphql/postgres"
	"github.com/stretchr/testify/require"
	"testing"
)

func TeardownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()
	err := db.Truncate(ctx)
	require.NoError(t, err)
}
