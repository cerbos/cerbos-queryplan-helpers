package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
	"entgo.io/ent/schema/edge"
)

// Company holds the schema definition for the Company entity.
type Company struct {
	ent.Schema
}

/*
model Company {
  id        String    @id @default(cuid())
  createdAt DateTime  @default(now())
  updatedAt DateTime  @updatedAt
  name      String
  contacts  Contact[]
}
*/
// Fields of the Company.
func (Company) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now),
		field.String("name"),
	}
}

// Edges of the Company.
func (Company) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("contacts", Contact.Type),
	}
}
