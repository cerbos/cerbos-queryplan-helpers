// Copyright 2021-2025 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/db"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/ent"
	"github.com/cerbos/cerbos-sdk-go/cerbos"
	enginev1 "github.com/cerbos/cerbos/api/genpb/cerbos/engine/v1"
	"github.com/ghodss/yaml"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed db/testdata/query_plans.yaml
var yamlBytes []byte

type Test struct {
	Input json.RawMessage `json:"input"`
	SQL   string          `json:"sql"`
	Args  []interface{}   `json:"args"`
}

func Test_BuildPredicate(t *testing.T) {
	is := require.New(t)
	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	is.NoError(err)
	var tests []Test
	err = json.Unmarshal(jsonBytes, &tests)
	is.NoError(err)
	for _, tt := range tests {
		t.Run(tt.SQL, func(t *testing.T) {
			is := require.New(t)
			e := new(enginev1.PlanResourcesFilter_Expression_Operand)
			err := protojson.Unmarshal(tt.Input, e)
			is.NoError(err)
			p, err := BuildPredicate(e.Node.(*enginev1.PlanResourcesFilter_Expression_Operand_Expression))
			is.NoError(err)
			p.SetDialect(dialect.Postgres)
			q, args := p.Query()
			is.Equal(tt.SQL, q)
			is.Equal(tt.Args, args)
		})
	}
}

func runCerbos(ctx context.Context, t *testing.T) string {
	t.Helper()

	is := require.New(t)
	pool, err := dockertest.NewPool("")
	is.NoError(err, "Could not connect to docker: %s", err)

	_, currFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("could not detect current file directory")
	}

	srcDir := filepath.Join(filepath.Dir(currFile), "cerbos")

	options := &dockertest.RunOptions{
		Repository: "ghcr.io/cerbos/cerbos",
		Tag:        "0.47.0",
		Cmd:        []string{"server", "--config=/config/conf.yaml"},
		WorkingDir: srcDir,
	}

	resource, err := pool.RunWithOptions(options, func(config *docker.HostConfig) {
		config.Mounts = []docker.HostMount{
			{
				Target: "/config",
				Source: filepath.Join(srcDir, "config"),
				Type:   "bind",
			},
			{
				Target: "/policies",
				Source: filepath.Join(srcDir, "policies"),
				Type:   "bind",
			},
		}
	})
	is.NoError(err, "Could not start resource: %s", err)

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Errorf("Failed to cleanup resources: %v", err)
		}
	})

	deadline, ok := t.Deadline()
	if !ok {
		deadline = time.Now().Add(5 * time.Minute)
	}

	ctx, cancelFunc := context.WithDeadline(ctx, deadline)
	t.Cleanup(cancelFunc)

	port := resource.GetPort("3592/tcp")
	cerbosAddr := fmt.Sprintf("127.0.0.1:%s", port)
	t.Log(cerbosAddr)
	healthEndpoint := fmt.Sprintf("http://%s/_cerbos/health", cerbosAddr)
	is.NoError(pool.Retry(func() error {
		if err := ctx.Err(); err != nil {
			return err
		}
		request, err := http.NewRequestWithContext(ctx, "GET", healthEndpoint, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errors.New("health check request status not OK")
		}
		return nil
	}), "Cerbos container did not start")

	return cerbosAddr
}

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	cerbosAddr := runCerbos(ctx, t)

	is := require.New(t)
	nick, simon, mary, christina, aleks := "Nick", "Simon", "Mary", "Christina", "Aleks"

	tests := []struct {
		username string
		want     []string
	}{
		{
			username: "alice",
			want:     []string{nick, simon, mary, christina, aleks},
		},
		{
			username: "john",
			want:     []string{nick, simon, mary, aleks},
		},
		{
			username: "sarah",
			want:     []string{mary, christina, aleks, nick},
		},
		{
			username: "geri",
			want:     []string{nick, aleks},
		},
	}

	c, err := cerbos.New(cerbosAddr, cerbos.WithPlaintext())
	is.NoError(err)
	repo, err := db.New(BuildPredicateType(BuildPredicate), ent.Log(t.Log), ent.Debug())
	is.NoError(err)
	err = repo.SetupDatabase(ctx)
	is.NoError(err)
	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			user, err := repo.GetUserByUsername(ctx, tt.username)
			is.NoError(err)
			// Create a new principal object with information from the database and the request.
			principal := cerbos.NewPrincipal(strconv.Itoa(user.ID)).
				WithRoles(user.Role).
				WithAttr("department", user.Department)

			queryPlan, err := c.PlanResources(ctx, principal, cerbos.NewResource("contact", ""), "read")
			is.NoError(err)

			filter := queryPlan.GetFilter()
			// t.Log(protojson.Format(filter))
			contacts, err := repo.GetContacts(ctx, filter)
			is.NoError(err)
			is.ElementsMatch(getNames(contacts), tt.want, tt.username)
		})
	}
}

func getNames(cs []*ent.Contact) []string {
	ns := make([]string, len(cs))
	for i := range cs {
		ns[i] = cs[i].FirstName
	}
	return ns
}
