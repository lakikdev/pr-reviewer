package userRole

import (
	"pr-reviewer/internal/model"

	"github.com/pkg/errors"
)

const revokeUserRoleQuery = `
	DELETE FROM user_roles
	WHERE user_id = $1 AND role = $2;
`

func (d *DB) Revoke(userID model.UserID, role model.Role) error {
	if _, err := d.tx.Exec(revokeUserRoleQuery, userID, role); err != nil {
		return errors.Wrap(err, "could not revoke user role")
	}
	return nil
}
