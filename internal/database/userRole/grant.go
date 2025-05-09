package userRole

import (
	"pr-reviewer/internal/model"

	"github.com/pkg/errors"
)

const grantUserRoleQuery = `
	INSERT INTO user_roles (user_id, role)
		VALUES ($1, $2);
`

func (d *DB) Grant(userID model.UserID, role model.Role) error {
	if _, err := d.tx.Exec(grantUserRoleQuery, userID, role); err != nil {
		return errors.Wrap(err, "could not grant user role")
	}
	return nil
}
