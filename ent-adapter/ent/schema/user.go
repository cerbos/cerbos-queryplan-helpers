package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

/*
model User {
  id         String    @id @default(cuid())
  username   String    @unique
  email      String    @unique
  name       String?
  contacts   Contact[]
  role       String
  department String
}

*/
// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Unique(),
		field.String("email").Unique(),
		field.String("name").Optional(),
		field.String("role"),
		field.String("department"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("contacts", Contact.Type),
	}
}
