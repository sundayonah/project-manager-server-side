package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Projects struct {
	ent.Schema
}

func (Projects) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("imageUrl").Optional(),
		field.String("link").Optional(),
		field.String("description").
			Optional().
			MaxLen(1000).
			Comment("A brief description of the package"),
		field.String("stacks").Default("[]"),
		// field.Time("created_at").
		// 	Default(time.Now).
		// 	Comment("The time the package was created"),
		// field.Time("updated_at").
		// 	Default(time.Now).
		// 	UpdateDefault(time.Now).
		// 	Comment("The time the package was last updated"),
	}
}
