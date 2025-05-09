package authHelper

import (
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
)

// we will have 3 permission type for now.
type PermissionType string

const (
	//User has 'admin' role
	Admin PermissionType = "admin"
	//User is logged in (we have userID in principal)
	Member PermissionType = "member"
	//User is logged in and user id passed to API is the same
	MemberIsTarget PermissionType = "memberIsTarget"
	//API-Key was provided and is valid
	APIKey PermissionType = "apiKey"
	//Any one can access
	Any PermissionType = "anonym"
)

type PermissionFunc func(c echo.Context, p *permissions, principal model.Principal, permissionType PermissionType) (bool, error)

var permissionsFuncsMap = map[PermissionType]PermissionFunc{
	Admin:          withRole,
	Member:         isLoggedIn,
	MemberIsTarget: memberIsTarget,
	APIKey:         apiKey,

	Any: any,
}

func any(c echo.Context, p *permissions, principal model.Principal, permissionType PermissionType) (bool, error) {
	return true, nil
}

func withRole(c echo.Context, p *permissions, principal model.Principal, permissionType PermissionType) (bool, error) {
	if hasRole, _ := p.withRole(principal, model.Role(permissionType)); hasRole {
		return true, nil
	}

	return false, nil
}

// Logged in user
func isLoggedIn(c echo.Context, p *permissions, principal model.Principal, permissionType PermissionType) (bool, error) {
	return member(principal), nil
}

// Logged in User = Target User
func memberIsTarget(c echo.Context, p *permissions, principal model.Principal, permissionType PermissionType) (bool, error) {
	targetUserID := model.UserID(c.Param("userID"))

	if targetUserID == model.NilUserID || principal.UserID == model.NilUserID {
		return false, nil
	}

	if targetUserID != principal.UserID {
		return false, nil
	}
	return true, nil
}

func apiKey(c echo.Context, p *permissions, principal model.Principal, permissionType PermissionType) (bool, error) {
	//get apiKey from request
	hashedAPIKey := GetAPIKeyFromContext(c.Request())
	if hashedAPIKey == "" {
		return false, nil
	}
	//check if apiKey is valid
	validAPIKey, err := p.dataRepository.APIKeysCache.Get(hashedAPIKey)
	if validAPIKey == nil || validAPIKey.(bool) == false {
		return false, err
	}

	return true, nil
}
