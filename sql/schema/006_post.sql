-- +goose Up
CREATE TABLE post (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    pub_date TIMESTAMPTZ,
    description TEXT NOT NULL,
    feed_id UUID NOT NULL REFERENCES feed(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE post;