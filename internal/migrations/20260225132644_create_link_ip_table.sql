-- +goose Up
CREATE TABLE link_ip (
    id BIGSERIAL PRIMARY KEY,
    link_id UUID NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    ip INET NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_link_ip_link_id ON link_ip(link_id);
CREATE INDEX idx_link_ip_ip ON link_ip(ip);

-- +goose Down
DROP TABLE link_ip;