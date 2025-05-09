package userAuth

import (
	"fmt"

	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/model"

	"github.com/pkg/errors"
)

const insertUserAuthQuery = `
	INSERT INTO user_auth (%s)
	VALUES (%s)
`

func (d *DB) AddToUser(userAuth *model.UserAuth) error {
	fields := dbHelper.GetDBFieldsWithIgnore(userAuth, []string{"auth_id", "active", "added_at"})
	insertQuery := fmt.Sprintf(insertUserAuthQuery, dbHelper.GetDBFieldsCSV(fields), dbHelper.GetDBFieldsCSVColons(fields))

	_, err := d.tx.NamedQuery(insertQuery, userAuth)
	if err != nil {
		return errors.Wrap(err, "could not add auth option for user")
	}

	return nil
}
