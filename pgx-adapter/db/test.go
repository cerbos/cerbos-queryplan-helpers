package db

import (
	"net"
	"strconv"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
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

func StartEmbeddedPostgres() (port int, cleanup func() error, err error) {
	port, err = GetFreePort()
	if err != nil {
		return 0, nil, err
	}

	pgConf := embeddedpostgres.DefaultConfig().Port(uint32(port))
	pg := embeddedpostgres.NewDatabase(pgConf)
	err = pg.Start()
	if err != nil {
		return 0, nil, err
	}

	return port, pg.Stop, nil
}
