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
    semester_activity_id UUID NOT NULL REFERENCES semester_activity(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    month month_enum NOT NULL,
    result BOOL,
    comment TEXT,
    UNIQUE(semester_activity_id, month, student_id)
);
-- +goose Down 
DROP TYPE month_enum;
DROP TABLE attestation;