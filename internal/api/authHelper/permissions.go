package authHelper

import (
	"net/http"

	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/common"
	"pr-reviewer/internal/model"

	"github.com/labstack/echo/v4"
)

// we probably need a better name
type PermissionsUserIDExtraction struct {
	UserID model.UserID `param:"userID"`
}

type Permissions interface {
	Wrap(next echo.HandlerFunc, permissionTypes ...PermissionType) echo.HandlerFunc
	Check(c echo.Context, permissionTypes ...PermissionType) bool
	RemoveFromCache(userID model.UserID) bool
	GetRoles(userID model.UserID) ([]model.Role, error)

	IsAdmin(userID model.UserID) (bool, error)
}

type permissions struct {
	dataRepository *common.DataRepository
}

func NewPermissions(dataRepository *common.DataRepository) Permissions {
	p := &permissions{
		dataRepository: dataRepository,
	}

	return p
}

func (p *permissions) RemoveFromCache(userID model.UserID) bool {
	return p.dataRepository.UserRolesCache.Remove(userID)
}

// get user's roles from cache (if we wont have roles in cache it will get it from database)
func (p *permissions) GetRoles(userID model.UserID) ([]model.Role, error) {
	roles, err := p.dataRepository.UserRolesCache.Get(userID)
	if err != nil {
		return nil, err
	}

	rolesList := make([]model.Role, 0)
	for _, role := range roles.([]*model.UserRole) {
		rolesList = append(rolesList, role.Role)
	}

	return rolesList, nil
}

func (p *permissions) withRole(principal model.Principal, role model.Role) (bool, error) {
	if principal.UserID == model.NilUserID {
		return false, nil
	}

	//Load roles
	roles, err := p.GetRoles(principal.UserID)
	if err != nil {
		return false, err
	}

	return hasRole(roles, role), nil
}

// We need to see if we have principal on Request in this point...
func (p *permissions) Wrap(next echo.HandlerFunc, permissionTypes ...PermissionType) echo.HandlerFunc {
	return func(c echo.Context) error {
		if allowed := p.Check(c, permissionTypes...); !allowed {
			//TODO: Update this error message to have more options, ex: if user fails MemberOwnsBundle check
			return utils.NewHttpError(nil, http.StatusUnauthorized, "Err401x9999", "Error while checking access token - Access Denied")
		}

		return next(c)
	}
}

func (p *permissions) IsAdmin(userID model.UserID) (bool, error) {
	return p.withRole(model.Principal{UserID: userID}, model.RoleAdmin)
}

// The idea is to return TRUE if one of permission types matches.
// for example if permission type is Admin and MemberIsTarget
// Admin can edit any user so if user has Admin role we don't care, admin don't match MemberIsTarget permission
func (p *permissions) Check(c echo.Context, permissionTypes ...PermissionType) bool {
	principal := GetPrincipal(c.Request())
	for _, permissionType := range permissionTypes {
		if permissionFunc, ok := permissionsFuncsMap[permissionType]; ok {
			if allowed, _ := permissionFunc(c, p, principal, permissionType); allowed {
				return true
			}
		}
	}
	return false
}
