package user

import (
	"pr-reviewer/internal/model"

	"github.com/pkg/errors"
)

const createUserQuery = `
	INSERT INTO users (user_id, email)
	VALUES (default, :email)
	RETURNING user_id;
`

func (d *DB) Create(user *model.User) error {

	rows, err := d.tx.NamedQuery(createUserQuery, user)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		return errors.Wrap(err, "could not create user")
	}

	rows.Next()
	if err := rows.Scan(&user.ID); err != nil {
		return errors.Wrap(err, "could not get created userID")
	}

	return nil
}
