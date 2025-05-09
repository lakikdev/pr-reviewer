CREATE TYPE user_auth_types AS ENUM (
    'udid',
    'credentials',
    'google',
    'apple',
    'facebook',
    'sequenceWallet'
);

CREATE TABLE user_auth (
    auth_id UUID DEFAULT uuid_generate_v4(), 
    user_id UUID NOT NULL REFERENCES users,
    type user_auth_types,
    data jsonb, 
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    active BOOLEAN NOT NULL DEFAULT 'true',
    PRIMARY KEY(user_id, type)
);

--Create index for auth types
CREATE INDEX user_auth_user
    ON user_auth (user_id)

--We need to create index for every user_auth_types to make sure we have fast logins
