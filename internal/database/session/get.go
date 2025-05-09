package session

import (
	"pr-reviewer/internal/model"
)

const getSessionQuery = `
	SELECT user_id, device_id, refresh_token_hash, expires_at, os_version, app_version, device_type, ip_address, last_updated
	FROM sessions
	WHERE user_id = $1
		AND device_id = $2
		AND to_timestamp(expires_at) > NOW()
`

func (d *DB) Get(data model.Session) (*model.Session, error) {
	var session model.Session
	if err := d.tx.Get(&session, getSessionQuery, data.UserID, data.DeviceID); err != nil {
		return nil, err
	}

	return &session, nil
}
