package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Positive(),
		field.String("name").Default("unknown"),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return nil
}
