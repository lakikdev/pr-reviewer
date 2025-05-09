package cacheUtils

import "errors"

// ErrNoData is returned when a record can't be found.
var ErrNoData = errors.New("item not found")
