package apiKey

import (
	"pr-reviewer/internal/model"

	"github.com/pkg/errors"
)

const createAPIKeyQuery = `
	INSERT INTO api_keys (name, key_hash, expires_at)
	VALUES (:name, :key_hash, :expires_at);
`

func (d *DB) Create(apiKey *model.APIKey) error {
	_, err := d.tx.NamedExec(createAPIKeyQuery, apiKey)
	if err != nil {
		return errors.Wrap(err, "could not create apiKey")
	}

	return nil
}
