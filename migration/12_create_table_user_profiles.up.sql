CREATE TABLE IF NOT EXISTS user_profiles(
    id SERIAL PRIMARY KEY,
    user_id SMALLINT,
    name VARCHAR(50),
    role SMALLINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);