package user

import (
	"fmt"

	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/model"
)

const updateUserQuery = `
	UPDATE users 
	SET %s
	WHERE user_id = :user_id and deleted_at is null;
`

func (d *DB) Update(user *model.User) error {
	fields := dbHelper.GetDBFieldsWithIgnore(user, []string{"user_id", "created_at", "is_creator", "stripe_customer_id", "stripe_payment_method_id"})
	updateQuery := fmt.Sprintf(updateUserQuery, dbHelper.GetDBFieldsUpdate(fields))

	result, err := d.tx.NamedExec(updateQuery, user)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return dbHelper.ErrItemNotFound
	}

	return nil
}
