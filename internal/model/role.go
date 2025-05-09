package model

// Role is a function a user can serve
type Role string

const (
	// RoleAdmin is an administrator of App. Root
	RoleAdmin     Role = "admin"
	RoleCreator   Role = "creator"
	RoleModerator Role = "moderator"
	RoleSelf      Role = "self"
)

//UserRole
type UserRole struct {
	UserID UserID `json:"userID" db:"user_id"`
	Role   Role   `json:"role" db:"role"`
}
