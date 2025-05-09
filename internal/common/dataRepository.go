package common

import (
	"context"
	"pr-reviewer/internal/database"
	"pr-reviewer/internal/model"
	"time"

	"github.com/bluele/gcache"
)

// Powerful set of commonly used data repositories.
type DataRepository struct {
	DB *database.DB

	UserRolesCache gcache.Cache
	APIKeysCache   gcache.Cache
}

func NewDataRepository(db *database.DB) *DataRepository {
	r := &DataRepository{
		DB: db,
	}

	r.UserRolesCache = gcache.New(200).
		LRU().
		LoaderExpireFunc(func(key interface{}) (interface{}, *time.Duration, error) {
			userID := key.(model.UserID)

			tx, _ := r.DB.BeginTxx(context.Background())
			roles, err := tx.UserRole().ListByUser(userID)
			if err != nil {
				_ = tx.Rollback()
				return nil, nil, err
			}
			_ = tx.Commit()

			expire := 1 * time.Minute
			return roles, &expire, nil
		}).
		Build()

	r.APIKeysCache = gcache.New(200).
		LRU().
		LoaderExpireFunc(func(key interface{}) (interface{}, *time.Duration, error) {
			hashedKey := key.(string)
			tx, _ := r.DB.BeginTxx(context.Background())
			apiKey, err := tx.APIKey().Exists(hashedKey)
			if err != nil {
				_ = tx.Rollback()
				return nil, nil, err
			}
			_ = tx.Commit()
			if !apiKey {
				return nil, nil, nil
			}
			expire := 1 * time.Minute
			return true, &expire, nil
		}).Build()

	return r
}
