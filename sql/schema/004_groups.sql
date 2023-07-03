-- +goose Up 
CREATE TABLE groups(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    subcode TEXT NOT NULL,
    stream UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    course smallint NOT NULL,
    UNIQUE(stream, subcode)
);
-- +goose Down
DROP TABLE groups;