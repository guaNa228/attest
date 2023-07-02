-- +goose Up 
CREATE TABLE streams (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    code TEXT UNIQUE NOT NULL,
    program UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    UNIQUE(name, program)
);
-- +goose Down 
DROP TABLE streams;