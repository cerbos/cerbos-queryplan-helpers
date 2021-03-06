// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CompaniesColumns holds the columns for the "companies" table.
	CompaniesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
	}
	// CompaniesTable holds the schema information for the "companies" table.
	CompaniesTable = &schema.Table{
		Name:       "companies",
		Columns:    CompaniesColumns,
		PrimaryKey: []*schema.Column{CompaniesColumns[0]},
	}
	// ContactsColumns holds the columns for the "contacts" table.
	ContactsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "first_name", Type: field.TypeString},
		{Name: "last_name", Type: field.TypeString},
		{Name: "active", Type: field.TypeBool},
		{Name: "marketing_opt_in", Type: field.TypeBool},
		{Name: "company_contacts", Type: field.TypeInt, Nullable: true},
		{Name: "user_contacts", Type: field.TypeInt, Nullable: true},
	}
	// ContactsTable holds the schema information for the "contacts" table.
	ContactsTable = &schema.Table{
		Name:       "contacts",
		Columns:    ContactsColumns,
		PrimaryKey: []*schema.Column{ContactsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "contacts_companies_contacts",
				Columns:    []*schema.Column{ContactsColumns[7]},
				RefColumns: []*schema.Column{CompaniesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "contacts_users_contacts",
				Columns:    []*schema.Column{ContactsColumns[8]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "username", Type: field.TypeString, Unique: true},
		{Name: "email", Type: field.TypeString, Unique: true},
		{Name: "name", Type: field.TypeString, Nullable: true},
		{Name: "role", Type: field.TypeString},
		{Name: "department", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CompaniesTable,
		ContactsTable,
		UsersTable,
	}
)

func init() {
	ContactsTable.ForeignKeys[0].RefTable = CompaniesTable
	ContactsTable.ForeignKeys[1].RefTable = UsersTable
}
