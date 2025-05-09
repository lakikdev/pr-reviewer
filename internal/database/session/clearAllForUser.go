package session

import (
	"pr-reviewer/internal/model"
)

const clearAllForUserQuery = `
	DELETE FROM sessions
	WHERE user_id = $1`

func (d *DB) ClearAllForUser(userID model.UserID) error {
	if _, err := d.tx.Exec(clearAllForUserQuery, userID); err != nil {
		return err
	}

	return nil
}
