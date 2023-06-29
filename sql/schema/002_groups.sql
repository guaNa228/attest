-- +goose Up 
CREATE TABLE groups (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    code TEXT UNIQUE NOT NULL
);
-- +goose Down 
DROP TABLE groups;