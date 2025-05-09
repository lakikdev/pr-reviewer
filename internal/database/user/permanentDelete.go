package user

import (
	"pr-reviewer/internal/model"
)

const permanentDeleteUserQuery = `
	DELETE FROM users WHERE user_id = $1;
`

func (d *DB) PermanentDelete(userID model.UserID) (bool, error) {
	result, err := d.tx.Exec(permanentDeleteUserQuery, userID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
