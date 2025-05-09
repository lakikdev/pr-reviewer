package apiKey

const apiKeyExistsQuery = `
	SELECT EXISTS (
		SELECT 1
		FROM api_keys
		WHERE key_hash = $1 and active = true and (expires_at IS NULL OR expires_at > now())
	)
`

func (d *DB) Exists(hashedKey string) (bool, error) {
	var exists bool

	err := d.tx.Get(&exists, apiKeyExistsQuery, hashedKey)
	if err != nil {
		return false, err
	}
	return exists, nil
}
