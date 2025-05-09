package dbHelper

import "errors"

//UniqueViolation Postgres error string for a unique index violation
const UniqueViolation = "unique_violation"

// ErrItemNotFound is returned when a record can't be found.
var ErrItemNotFound = errors.New("item not found")
var ErrItemExists = errors.New("item already exists")

// IsUniqueViolation checks if the error is a unique violation
func IsUniqueViolation(err error) bool {
	return err != nil && err.Error() == UniqueViolation
}
