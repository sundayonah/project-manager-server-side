package schema

import (
	"time"

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
			NotEmpty().
			Unique().
			MaxLen(100).
			Comment("The name of the package"),
		field.String("link").
			Optional().
			Comment("The link to the package"),
		field.String("description").
			Optional().
			MaxLen(1000).
			Comment("A brief description of the package"),
		field.String("stacks").Default("[]"),
		field.Time("created_at").
			Default(time.Now).
			Comment("The time the package was created"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The time the package was last updated"),
	}
}
