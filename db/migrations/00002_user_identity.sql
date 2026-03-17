-- +goose Up
CREATE TABLE user_identities (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider     VARCHAR(50) NOT NULL DEFAULT 'local',
    password     TEXT,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE user_identities;