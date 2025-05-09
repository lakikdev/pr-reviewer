CREATE TABLE sessions (
    user_id UUID REFERENCES users,
    device_id TEXT,
    refresh_token_hash bytea,
    expires_at INTEGER,
    ip_address TEXT,
    os_version TEXT,
    app_version TEXT,
    device_type TEXT,
    last_updated TIMESTAMP,
    PRIMARY KEY (user_id, device_id)
);
