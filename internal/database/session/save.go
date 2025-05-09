package session

import (
	"pr-reviewer/internal/model"
)

const insertOrUpdateSession = `
	INSERT INTO sessions (user_id, device_id, refresh_token_hash, expires_at, ip_address, os_version, app_version, device_type, last_updated)
	VALUES(:user_id, :device_id, :refresh_token_hash, :expires_at, :ip_address, :os_version, :app_version, :device_type, NOW())

	ON CONFLICT (user_id, device_id)
	DO
		UPDATE
			SET refresh_token_hash = :refresh_token_hash,
				expires_at = :expires_at,
				ip_address = :ip_address,
				os_version = :os_version, 
				app_version = :app_version, 
				device_type = :device_type,
				last_updated = NOW();
`

func (d *DB) Save(session model.Session) error {
	if _, err := d.tx.NamedExec(insertOrUpdateSession, session); err != nil {
		return err
	}

	return nil
}
