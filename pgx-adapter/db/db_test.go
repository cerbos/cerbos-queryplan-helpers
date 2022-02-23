// Copyright 2021-2022 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func Test_SetupDatabase(t *testing.T) {
	is := require.New(t)

	port, cleanup, err := StartEmbeddedPostgres()
	is.NoError(err)

	t.Cleanup(func() {
		if err := cleanup(); err != nil {
			t.Errorf("Failed to stop Postgres: %v", err)
		}
	})

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)

	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://postgres:postgres@localhost:%d?sslmode=disable", port))
	is.NoError(err)
	t.Cleanup(func() {
		conn.Close(ctx)
	})
	err = SetupDatabase(ctx, conn)
	is.NoError(err)

	connURL := fmt.Sprintf("postgres://cerbforce_user:cerb@localhost:%d/postgres?sslmode=disable&search_path=cerbforce", port)
	repo, err := New(ctx, nil, connURL)
	is.NoError(err)

	t.Cleanup(func() {
		repo.client.Close(ctx)
	})

	user, err := repo.GetUserByUsername(ctx, "alice")
	is.NoError(err)
	is.Equal(user.Department, "IT")
	user, err = repo.GetUserByUsername(ctx, "no-such-username")
	is.NoError(err)
	is.Nil(user)
}
