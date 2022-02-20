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
