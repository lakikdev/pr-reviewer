package authHelper

import (
	"pr-reviewer/internal/model"
)

var hasRole = func(roles []model.Role, userRole model.Role) bool {
	for _, role := range roles {
		if role == userRole {
			return true
		}
	}
	return false
}

var member = func(principal model.Principal) bool {
	return principal.UserID != model.NilUserID
}
