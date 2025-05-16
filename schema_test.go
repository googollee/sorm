package sorm

import (
	"testing"
	"time"
)

type TestUser struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	PasswordSalt string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DisabledAt   time.Time
}

func (u *TestUser) Schema(schema *Schema) {
	schema.Field(&u.ID).Name("id").PrimaryKey()
	schema.Field(&u.Name).Name("name").NotNull().Type(VarChar(64)).UniqueIndex()
	schema.Field(&u.Email).Name("email").NotNull().Type(VarChar(128))
	schema.Field(&u.PasswordHash).Name("password_hash").NotNull().Type(VarChar(128))
	schema.Field(&u.PasswordSalt).Name("password_salt").NotNull().Type(VarChar(128))
	schema.Field(&u.CreatedAt).Name("created_at").NotNull().Type(IntegerByte(8))
	schema.Field(&u.UpdatedAt).Name("updated_at").NotNull().Type(IntegerByte(8))
	schema.Field(&u.DisabledAt).Name("disabled_at").Nullable(true).Type(IntegerByte(8))
	schema.UniqueIndex(&u.Name, &u.Email)
}

func TestSchemaSimple(t *testing.T) {
	var u TestUser
	schema := newSchema(&u)
	u.Schema(schema)
}
