-- +goose Up 
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    teacher_id integer UNIQUE
);
-- +goose StatementBegin
CREATE FUNCTION validate_teacher_id() RETURNS TRIGGER AS $BODY$ BEGIN IF NEW.teacher_id IS NOT NULL
AND NEW.role <> 'teacher' THEN RAISE EXCEPTION 'If teacher_id is not null, role must be ''teacher''';
END IF;
RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER validate_teacher_id_trigger BEFORE
INSERT
    OR
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION validate_teacher_id();
-- +goose Down 
DROP TRIGGER validate_teacher_id_trigger on users;
DROP FUNCTION validate_teacher_id;
DROP TABLE users;