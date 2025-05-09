package user

import (
	"database/sql"
	"fmt"

	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/model"
)

const getUserByIDQuery = `
	SELECT %s
	FROM users
	WHERE user_id = $1 and deleted_at is null;
`

func (d *DB) GetByID(userID model.UserID) (*model.User, error) {
	var user model.User
	fields := dbHelper.GetDBFields(user)
	selectQuery := fmt.Sprintf(getUserByIDQuery, dbHelper.GetDBFieldsCSV(fields))

	if err := d.tx.Get(&user, selectQuery, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, dbHelper.ErrItemNotFound
		}
		return nil, err
	}

	return &user, nil
}
