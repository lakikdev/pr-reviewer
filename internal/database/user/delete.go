package user

import (
	"pr-reviewer/internal/model"
)

const deleteUserQuery = `
	UPDATE users
	SET deleted_at = NOW()
	WHERE user_id = $1 AND deleted_at IS NULL;
`

const deleteAllAuthForUser = `
	DELETE FROM user_auth WHERE user_id = $1;
`

func (d *DB) Delete(userID model.UserID) error {
	_, err := d.tx.Exec(deleteUserQuery, userID)
	if err != nil {
		return err
	}

	// Delete all auth for user to make sure we can reuse them
	_, err = d.tx.Exec(deleteAllAuthForUser, userID)
	if err != nil {
		return err
	}

	return nil
}
