-- +goose Up
SELECT 'up SQL query';
CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL UNIQUE,
    long_url TEXT NOT NULL,
    click_count BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_links_user_id ON links(user_id);
-- +goose Down
SELECT 'down SQL query';
DROP TABLE links;