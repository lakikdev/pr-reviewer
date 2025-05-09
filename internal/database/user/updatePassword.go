package user

import (
	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/model"
)

const updatePasswordQuery = `
	UPDATE users 
	SET password_hash = :password_hash
	WHERE user_id = :user_id;
`

func (d *DB) UpdatePassword(user *model.User) error {
	result, err := d.tx.NamedExec(updatePasswordQuery, user)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return dbHelper.ErrItemNotFound
	}

	return nil
}
