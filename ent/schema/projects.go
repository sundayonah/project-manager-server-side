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
		field.String("image_url").Optional(),
		field.String("link").Optional(),
		field.String("description").Optional(),
	}
}
