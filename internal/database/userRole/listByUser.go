package userRole

import (
	"pr-reviewer/internal/model"

	"github.com/pkg/errors"
)

const getRolesByUserIDQuery = `
	SELECT role
	FROM user_roles
	WHERE user_id = $1;
`

func (d *DB) ListByUser(userID model.UserID) ([]*model.UserRole, error) {
	var roles []*model.UserRole
	if err := d.tx.Select(&roles, getRolesByUserIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "could not get user roles")
	}

	return roles, nil
}
