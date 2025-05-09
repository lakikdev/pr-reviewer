package session

import (
	"pr-reviewer/internal/model"
)

const getLatestSessionQuery = `
	SELECT user_id, device_id, refresh_token_hash, expires_at,  COALESCE(os_version, 'N/A') as os_version, 
		COALESCE(app_version, 'N/A') as app_version, COALESCE(device_type, 'N/A') as device_type, 
		COALESCE(ip_address, 'N/A') as ip_address, last_updated
	FROM sessions
	WHERE user_id = $1
	ORDER BY expires_at desc 
	LIMIT 1;
`

func (d *DB) GetLatest(userID model.UserID) (*model.Session, error) {
	var session model.Session
	if err := d.tx.Get(&session, getLatestSessionQuery, userID); err != nil {
		return nil, err
	}

	return &session, nil
}
