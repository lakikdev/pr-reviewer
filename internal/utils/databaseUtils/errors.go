package databaseUtils

import "errors"

//UniqueViolation Postgres error string for a unique index violation
const UniqueViolation = "unique_violation"

// ErrItemNotFound is returned when a record can't be found.
var ErrItemNotFound = errors.New("item not found")

var ErrItemAlreadyExists = errors.New("item with that id already exists")
