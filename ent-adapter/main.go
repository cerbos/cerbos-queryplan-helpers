package main

import (
	"net/http"
	"log"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/ent"
	"context"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/service"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/db"
	"flag"
)

func main() {
	http.Handle("/", initHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initHandler() *service.Service {
	cerbosAddr := flag.String("cerbos", "localhost:3592", "Address of the Cerbos server")
	flag.Parse()

	client, err := db.New(ent.Log(log.Println), ent.Debug())
	if err != nil {
		log.Fatalf("Could not create a database client: %s", err)
	}
	ctx := context.Background()
	client.SetupDatabase(ctx)
	if err != nil {
		log.Fatalf("Could not set up a database: %s", err)
	}

	s, err := service.NewService(*cerbosAddr, client)
	if err != nil {
		log.Fatalf("Could not create a service: %s", err)
	}
	return s
}
