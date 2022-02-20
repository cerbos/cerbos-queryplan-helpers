package db

import (
	"context"
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/ent/predicate"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/ent"
	"encoding/json"
	"entgo.io/ent/dialect/sql"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/ent/contact"
)

func Test_ReadSeedFile(t *testing.T) {
	is := require.New(t)
	// Create an ent.Client with in-memory SQLite database.
	c, err := New(ent.Log(t.Log), ent.Debug())
	client := c.client
	defer client.Close()
	ctx := context.Background()
	err = c.SetupDatabase(ctx)
	is.NoError(err)
	got := client.User.Query().CountX(ctx)
	var data []json.RawMessage
	err = json.Unmarshal(seed, &data)
	is.NoError(err)
	is.Equal(len(data), got)
	contact := client.Contact.
		Query().
		WithCompany().
		Where(maryJane()).
		OnlyX(ctx)
	company := contact.Edges.Company
	is.Equal("Pepsi Co", company.Name)
}

func maryJane() predicate.Contact {
	return func(s *sql.Selector) {
		s.Where(sql.And(
			sql.EQ(contact.FieldFirstName, "Mary"),
			sql.EQ(contact.FieldLastName, "Jane"),
		))
	}
	//return contact.And(
	//	contact.FirstName("Mary"),
	//	contact.LastName("Jane"))
}
