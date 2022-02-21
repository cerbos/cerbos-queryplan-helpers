package db

import (
	"testing"
	_ "embed"
	"github.com/stretchr/testify/require"
	"github.com/ghodss/yaml"
	"encoding/json"
	responsev1 "github.com/cerbos/cerbos/api/genpb/cerbos/response/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"entgo.io/ent/dialect"
	cerbos "github.com/cerbos/cerbos/client"
	"context"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/ent"
	"strconv"
)

//go:embed testdata/query_plans.yaml
var yamlBytes []byte

type Test struct {
	Input json.RawMessage `json:"input"`
	Sql   string          `json:"sql"`
	Args  []interface{}   `json:"args"`
}

func Test_getPredicate(t *testing.T) {
	is := require.New(t)
	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	is.NoError(err)
	var tests []Test
	err = json.Unmarshal(jsonBytes, &tests)
	is.NoError(err)
	for _, tt := range tests {
		t.Run(tt.Sql, func(t *testing.T) {
			is := require.New(t)
			e := new(responsev1.ResourcesQueryPlanResponse_Expression_Operand)
			err := protojson.Unmarshal(tt.Input, e)
			is.NoError(err)
			p, err := getPredicate(e.Node.(*responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression))
			is.NoError(err)
			p.SetDialect(dialect.Postgres)
			q, args := p.Query()
			is.Equal(tt.Sql, q)
			is.Equal(tt.Args, args)
		})
	}
}

func TestIntegration(t *testing.T) {
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
	cerbosAddr := "localhost:3592"
	c, err := cerbos.New(cerbosAddr, cerbos.WithPlaintext())
	is.NoError(err)
	ctx := context.Background()
	repo, err := New(ent.Log(t.Log), ent.Debug())
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

			queryPlan, err := c.ResourcesQueryPlan(ctx, principal, cerbos.NewResource("contact", ""), "read")
			is.NoError(err)

			filter := queryPlan.GetFilter()
			//t.Log(protojson.Format(filter))
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
