package userAuth

import (
	"database/sql"
	"fmt"

	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/model"
)

const getByDataQuery = `
	SELECT %s FROM user_auth
	where type = $1 and data ->> $2 = LOWER($3);
`

func (d *DB) GetByData(authType string, fieldName string, fieldValue interface{}) (*model.UserAuth, error) {
	var userAuth model.UserAuth
	fields := dbHelper.GetDBFields(userAuth)
	selectQuery := fmt.Sprintf(getByDataQuery, dbHelper.GetDBFieldsCSV(fields))

	if err := d.tx.Get(&userAuth, selectQuery, authType, fieldName, fieldValue); err != nil {
		if err == sql.ErrNoRows {
			return nil, dbHelper.ErrItemNotFound
		}
		return nil, err
	}

	return &userAuth, nil
}
