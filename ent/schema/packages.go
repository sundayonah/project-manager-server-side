package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Packages holds the schema definition for the Packages entity.
type Packages struct {
	ent.Schema
}

// Fields of the Packages.
func (Packages) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("link").
			Optional(),
		field.String("description").
			Optional(),
	}
}
