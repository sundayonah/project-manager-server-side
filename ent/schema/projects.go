package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Projects holds the schema definition for the Projects entity.
type Projects struct {
	ent.Schema
}

// Fields of the Projects.
func (Projects) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("image_url").
			Optional(),
		field.String("link").
			Optional(),
		field.String("description").
			Optional(),
	}
}
