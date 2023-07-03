// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MonthEnum string

const (
	MonthEnumValue0  MonthEnum = "Январь"
	MonthEnumValue1  MonthEnum = "Февраль"
	MonthEnumValue2  MonthEnum = "Март"
	MonthEnumValue3  MonthEnum = "Апрель"
	MonthEnumValue4  MonthEnum = "Май"
	MonthEnumValue5  MonthEnum = "Июнь"
	MonthEnumValue6  MonthEnum = "Июль"
	MonthEnumValue7  MonthEnum = "Август"
	MonthEnumValue8  MonthEnum = "Сентябрь"
	MonthEnumValue9  MonthEnum = "Октябрь"
	MonthEnumValue10 MonthEnum = "Ноябрь"
	MonthEnumValue11 MonthEnum = "Декабрь"
)

func (e *MonthEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MonthEnum(s)
	case string:
		*e = MonthEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for MonthEnum: %T", src)
	}
	return nil
}

type NullMonthEnum struct {
	MonthEnum MonthEnum
	Valid     bool // Valid is true if MonthEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMonthEnum) Scan(value interface{}) error {
	if value == nil {
		ns.MonthEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MonthEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMonthEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MonthEnum), nil
}

type Attestation struct {
	ID                 uuid.UUID
	SemesterActivityID uuid.UUID
	StudentID          uuid.UUID
	Month              MonthEnum
	Result             sql.NullBool
	Comment            sql.NullString
}

type Program struct {
	ID         uuid.UUID `json:"id"`
	Name       string `json:"name"`
	MaxCourses int16 `json:"max_courses"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	Password  string `json:"-"`
	Role      string `json:"role"`
	TeacherID sql.NullInt32 `json:"teacher_id"`
	GroupID   uuid.NullUUID `json:"group_id"`
}

type Group struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
	Subcode   string `json:"subcode"`
	Stream    uuid.UUID `json:"stream"`
	Course    int16 `json:"course"`
}

type Stream struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Program   uuid.UUID `json:"program"`
}

//Old 
type Class struct {
	ID   uuid.UUID `json:"id"`
	Name string `json:"name"`
}

//Old 
type SemesterActivity struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ClassID   uuid.UUID `json:"class_id"`
	GroupID   uuid.UUID `json:"group_id"`
	TeacherID uuid.UUID `json:"teacher_id"`
}