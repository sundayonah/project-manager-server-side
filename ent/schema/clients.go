package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Clients holds the schema definition for the Clients entity.
type Clients struct {
	ent.Schema
}

// Fields of the Clients.
func (Clients) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique().
			MaxLen(100).
			Comment("The name of the package"),
		field.String("link").
			Optional().
			Comment("The link to the package"),
		field.String("imageUrl").Optional().Comment("The image URL of the client"),
		field.Time("created_at").
			Default(time.Now).
			Comment("The time the package was created"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The time the package was last updated"),
	}
}
