// Code generated by entc, DO NOT EDIT.

package company

import (
	"time"
)

const (
	// Label holds the string label denoting the company type in the database.
	Label = "company"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgeContacts holds the string denoting the contacts edge name in mutations.
	EdgeContacts = "contacts"
	// Table holds the table name of the company in the database.
	Table = "companies"
	// ContactsTable is the table that holds the contacts relation/edge.
	ContactsTable = "contacts"
	// ContactsInverseTable is the table name for the Contact entity.
	// It exists in this package in order to avoid circular dependency with the "contact" package.
	ContactsInverseTable = "contacts"
	// ContactsColumn is the table column denoting the contacts relation/edge.
	ContactsColumn = "company_contacts"
)

// Columns holds all SQL columns for company fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldName,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
)
