package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

func (X509Certificate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("x509_certificates").Unique(),
	}
}
