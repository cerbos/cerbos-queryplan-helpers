package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"net"
	"strconv"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"fmt"
	"github.com/jackc/pgx/v4"
)

func GetFreeListenAddr() (string, error) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	addr := lis.Addr().String()

	return addr, lis.Close()
}

func GetFreePort() (int, error) {
	addr, err := GetFreeListenAddr()
	if err != nil {
		return 0, err
	}

	_, p, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(p)
}

func Test_ReadSeedFile(t *testing.T) {
	is := require.New(t)
	port, err := GetFreePort()
	is.NoError(err, "Failed to get free port")

	pgConf := embeddedpostgres.DefaultConfig().Port(uint32(port))
	pg := embeddedpostgres.NewDatabase(pgConf)
	require.NoError(t, pg.Start(), "Failed to start Postgres")

	t.Cleanup(func() {
		if err := pg.Stop(); err != nil {
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
	c, err := New(ctx, nil, connURL)
	client := c.client
	t.Cleanup(func() {
		client.Close(ctx)
	})
	user, err := c.GetUserByUsername(ctx, "alice")
	is.NoError(err)
	is.Equal(user.Name, "Alice")
	user, err = c.GetUserByUsername(ctx, "no-such-username")
	is.NoError(err)
	is.Nil(user)
}
