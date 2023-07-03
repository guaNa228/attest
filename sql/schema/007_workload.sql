-- +goose Up
CREATE TABLE workload(
    id UUID PRIMARY KEY,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    class UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    teacher UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(group_id, class) --Убрать нахуй, так как у одной группы по одному предмету может быть несколько преподавателей,
    --которые должны иметь возможность проставлять аттестацию
    --ПОДУМОЙ
);
-- +goose StatementBegin
CREATE FUNCTION validate_teacher_role() RETURNS trigger LANGUAGE 'plpgsql' NOT LEAKPROOF AS $BODY$ BEGIN IF NOT EXISTS (
    SELECT 1
    FROM users
    WHERE id = NEW.teacher
        AND role = 'teacher'
) THEN RAISE EXCEPTION 'Referenced user is not a teacher' USING ERRCODE = '23503';
END IF;
RETURN NEW;
END;
$BODY$;
-- +goose StatementEnd
CREATE TRIGGER validate_teacher_role_trigger BEFORE
INSERT
    OR
UPDATE ON workload FOR EACH ROW EXECUTE FUNCTION validate_teacher_role();
-- +goose Down
DROP TRIGGER validate_teacher_role_trigger ON workload;
DROP FUNCTION validate_teacher_role;
DROP TABLE workload;