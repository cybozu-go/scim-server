// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/cybozu-go/scim-server/ent/email"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/google/uuid"
)

// Email is the model entity for the Email schema.
type Email struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Display holds the value of the "display" field.
	Display string `json:"display,omitempty"`
	// Primary holds the value of the "primary" field.
	Primary bool `json:"primary,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Value holds the value of the "value" field.
	Value string `json:"value,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EmailQuery when eager-loading is set.
	Edges       EmailEdges `json:"edges"`
	user_emails *uuid.UUID
}

// EmailEdges holds the relations/edges for other nodes in the graph.
type EmailEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EmailEdges) UserOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.User == nil {
			// The edge user was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.User, nil
	}
	return nil, &NotLoadedError{edge: "user"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Email) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case email.FieldPrimary:
			values[i] = new(sql.NullBool)
		case email.FieldDisplay, email.FieldType, email.FieldValue:
			values[i] = new(sql.NullString)
		case email.FieldID:
			values[i] = new(uuid.UUID)
		case email.ForeignKeys[0]: // user_emails
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			return nil, fmt.Errorf("unexpected column %q for type Email", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Email fields.
func (e *Email) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case email.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				e.ID = *value
			}
		case email.FieldDisplay:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field display", values[i])
			} else if value.Valid {
				e.Display = value.String
			}
		case email.FieldPrimary:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field primary", values[i])
			} else if value.Valid {
				e.Primary = value.Bool
			}
		case email.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				e.Type = value.String
			}
		case email.FieldValue:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field value", values[i])
			} else if value.Valid {
				e.Value = value.String
			}
		case email.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field user_emails", values[i])
			} else if value.Valid {
				e.user_emails = new(uuid.UUID)
				*e.user_emails = *value.S.(*uuid.UUID)
			}
		}
	}
	return nil
}

// QueryUser queries the "user" edge of the Email entity.
func (e *Email) QueryUser() *UserQuery {
	return (&EmailClient{config: e.config}).QueryUser(e)
}

// Update returns a builder for updating this Email.
// Note that you need to call Email.Unwrap() before calling this method if this Email
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Email) Update() *EmailUpdateOne {
	return (&EmailClient{config: e.config}).UpdateOne(e)
}

// Unwrap unwraps the Email entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Email) Unwrap() *Email {
	_tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Email is not a transactional entity")
	}
	e.config.driver = _tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Email) String() string {
	var builder strings.Builder
	builder.WriteString("Email(")
	builder.WriteString(fmt.Sprintf("id=%v, ", e.ID))
	builder.WriteString("display=")
	builder.WriteString(e.Display)
	builder.WriteString(", ")
	builder.WriteString("primary=")
	builder.WriteString(fmt.Sprintf("%v", e.Primary))
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(e.Type)
	builder.WriteString(", ")
	builder.WriteString("value=")
	builder.WriteString(e.Value)
	builder.WriteByte(')')
	return builder.String()
}

// Emails is a parsable slice of Email.
type Emails []*Email

func (e Emails) config(cfg config) {
	for _i := range e {
		e[_i].config = cfg
	}
}
