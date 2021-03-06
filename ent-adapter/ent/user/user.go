// Code generated by entc, DO NOT EDIT.

package user

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldUsername holds the string denoting the username field in the database.
	FieldUsername = "username"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldRole holds the string denoting the role field in the database.
	FieldRole = "role"
	// FieldDepartment holds the string denoting the department field in the database.
	FieldDepartment = "department"
	// EdgeContacts holds the string denoting the contacts edge name in mutations.
	EdgeContacts = "contacts"
	// Table holds the table name of the user in the database.
	Table = "users"
	// ContactsTable is the table that holds the contacts relation/edge.
	ContactsTable = "contacts"
	// ContactsInverseTable is the table name for the Contact entity.
	// It exists in this package in order to avoid circular dependency with the "contact" package.
	ContactsInverseTable = "contacts"
	// ContactsColumn is the table column denoting the contacts relation/edge.
	ContactsColumn = "user_contacts"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldUsername,
	FieldEmail,
	FieldName,
	FieldRole,
	FieldDepartment,
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
