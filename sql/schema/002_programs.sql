-- +goose Up 
CREATE TABLE programs (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    max_courses smallint NOT NULL
);
-- +goose Down 
DROP TABLE programs;