package user

import (
	"database/sql"

	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/model"
	"pr-reviewer/internal/utils/errorUtils"

	"github.com/pkg/errors"
)

func (d *DB) List(param model.ListDataParameters) (users []*model.User, total *int64, err error) {
	fields := dbHelper.GetDBFields(users)

	queryOptions := dbHelper.QueryOptions{
		TableName:          "users",
		TableColumns:       fields,
		DefaultWhereCases:  []string{"deleted_at is null"},
		DefaultSort:        []string{"created_at desc"},
		QuickFilterColumns: []string{"user_id", "display_name"},
		TargetObject:       users,
	}

	sqlQuery, args, err := dbHelper.CreateQuery(param, queryOptions)
	if err != nil {
		return nil, nil, errorUtils.Wrap(err)
	}

	if err := d.tx.Select(&users, sqlQuery, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, errors.Wrap(err, "could not get users")
	}

	queryOptions.TableColumns = []string{"count(*)"}

	sqlQuery, args, err = dbHelper.CreateQuery(param, queryOptions)
	if err != nil {
		return nil, nil, errorUtils.Wrap(err)
	}

	if err := d.tx.Get(&total, sqlQuery, args...); err != nil {
		return nil, nil, errors.Wrap(err, "could not get users total")
	}

	return users, total, err
}
