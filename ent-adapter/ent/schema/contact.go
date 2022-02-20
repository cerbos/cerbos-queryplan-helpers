package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
	"entgo.io/ent/schema/edge"
)

// Contact holds the schema definition for the Contact entity.
type Contact struct {
	ent.Schema
}

/*

model Contact {
  id             String   @id @default(cuid())
  createdAt      DateTime @default(now())
  updatedAt      DateTime @updatedAt
  firstName      String
  lastName       String
  owner          User     @relation(fields: [ownerId], references: [id])
  ownerId        String
  company        Company? @relation(fields: [companyId], references: [id])
  companyId      String?
  active         Boolean  @default(false)
  marketingOptIn Boolean  @default(false)
}
*/
// Fields of the Contact.
func (Contact) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now),
		field.String("first_name"),
		field.String("last_name"),
		field.Bool("active"),
		field.Bool("marketing_opt_in"),
	}
}

// Edges of the Contact.
func (Contact) Edges() []ent.Edge {
	return []ent.Edge{
		edge.
			From("company", Company.Type).
			Ref("contacts").
			Unique(),

		edge.
			From("owner", User.Type).
			Ref("contacts").
			Unique().
			Required(),
	}
}
