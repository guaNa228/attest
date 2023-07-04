-- +goose Up 
CREATE TYPE month_enum AS ENUM (
    'Январь',
    'Февраль',
    'Март',
    'Апрель',
    'Май',
    'Июнь',
    'Июль',
    'Август',
    'Сентябрь',
    'Октябрь',
    'Ноябрь',
    'Декабрь'
);
CREATE TABLE attestation (
    id UUID PRIMARY KEY,
    workload UUID NOT NULL REFERENCES workloads(id) ON DELETE CASCADE,
    student UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    month month_enum NOT NULL,
    result BOOL,
    comment TEXT,
    UNIQUE(workload, month, student)
);
-- +goose StatementBegin
CREATE FUNCTION validate_student_role() RETURNS trigger LANGUAGE 'plpgsql' NOT LEAKPROOF AS $BODY$ BEGIN IF NOT EXISTS (
    SELECT 1
    FROM users
    WHERE id = NEW.student
        AND role = 'student'
) THEN RAISE EXCEPTION 'Referenced user is not a student' USING ERRCODE = '23503';
END IF;
RETURN NEW;
END;
$BODY$;
CREATE TRIGGER validate_student_role_trigger BEFORE
INSERT
    OR
UPDATE ON attestation FOR EACH ROW EXECUTE FUNCTION validate_student_role();
-- +goose StatementEnd
-- +goose Down 
DROP TRIGGER validate_student_role_trigger on attestation;
DROP FUNCTION validate_student_role;
DROP TABLE attestation;
DROP TYPE month_enum;