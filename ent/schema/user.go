package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(`addresses`, Address.Type),
		edge.To(`groups`, GroupMember.Type),
		edge.To(`emails`, Email.Type),
		edge.To(`name`, Names.Type).
			Unique(),
		edge.To(`entitlements`, Entitlement.Type),
		edge.To(`roles`, Role.Type),
		edge.To(`IMS`, IMS.Type),
		edge.To(`phone_numbers`, PhoneNumber.Type),
		edge.To(`photos`, Photo.Type),
		edge.To(`x509_certificates`, X509Certificate.Type),
	}
}
