package datatype

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	mapset "github.com/deckarep/golang-set/v2"
)

var (
	TimeZone = "Asia/Shanghai"
)

const (
	DateTimeTZFormat = time.DateTime + "-07:00"
)

func LocalTimeZone() *time.Location {
	if v, err := time.LoadLocation(TimeZone); err == nil {
		return v
	}

	return time.UTC
}

func ParseDateTime(layout string, value string) (*time.Time, error) {
	value = strings.TrimSpace(value)

	if len(value) > len(layout) {
		value = value[:len(layout)]
	}

	value = strings.ReplaceAll(value, "T", " ")

	v, err := time.ParseInLocation(layout, value, LocalTimeZone())

	if err != nil {
		return nil, err
	}

	if v.Year() == 0 {
		t := time.Now().In(LocalTimeZone())

		v = time.Date(
			t.Year(), t.Month(), t.Day(), v.Hour(), v.Minute(), v.Second(), 0, v.Location(),
		)
	}

	return &v, nil
}

// ---------------------------------------------------------
//
//  DateTime
//
// ---------------------------------------------------------

type DateTime time.Time

func NewDateTime(t *time.Time) DateTime {
	v := time.Now()

	if t != nil {
		v = *t
	}

	v = v.In(LocalTimeZone())

	return DateTime(
		time.Date(
			v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second(), 0, v.Location(),
		),
	)
}

// GORM
func (dt DateTime) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "TIMESTAMPTZ"
	case "sqlite":
		return "TEXT"
	case "sqlserver":
		return "DATETIMEOFFSET"
	default:
		return "DATETIME"
	}
}

func (dt DateTime) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	v := time.Time(dt).In(LocalTimeZone())

	if mapset.NewSet("postgres", "sqlite", "sqlserver").Contains(db.Dialector.Name()) {
		return gorm.Expr("?", v.Format(DateTimeTZFormat))
	} else {
		return gorm.Expr("?", v.Format(time.DateTime))
	}
}

func (dt *DateTime) Scan(value any) error {
	layouts := []string{DateTimeTZFormat, time.DateTime}

	switch val := value.(type) {
	case []byte:
		for _, layout := range layouts {
			if v, err := ParseDateTime(layout, string(val)); err == nil {
				*dt = NewDateTime(v)

				break
			}
		}
	case string:
		for _, layout := range layouts {
			if v, err := ParseDateTime(layout, val); err == nil {
				*dt = NewDateTime(v)

				break
			}
		}
	case time.Time:
		*dt = NewDateTime(&val)
	}

	return nil
}

// JSON
func (dt DateTime) MarshalJSON() ([]byte, error) {
	v := time.Time(dt).In(LocalTimeZone())

	return json.Marshal(v.Format(time.DateTime))
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	v, err := ParseDateTime(time.DateTime, s)

	if err != nil {
		return err
	}

	*dt = NewDateTime(v)

	return nil
}

// String
func (dt DateTime) String() string {
	v := time.Time(dt).In(LocalTimeZone())

	return v.Format(time.DateTime)
}

// ---------------------------------------------------------
//
//  Date
//
// ---------------------------------------------------------

type Date time.Time

func NewDate(t *time.Time) Date {
	v := time.Now()

	if t != nil {
		v = *t
	}

	v = v.In(LocalTimeZone())

	return Date(
		time.Date(
			v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location(),
		),
	)
}

// GORM
func (d Date) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	if db.Dialector.Name() == "sqlite" {
		return "TEXT"
	} else {
		return "DATE"
	}
}

func (d *Date) Scan(value any) error {
	switch val := value.(type) {
	case []byte:
		if v, err := ParseDateTime(time.DateOnly, string(val)); err == nil {
			*d = NewDate(v)
		}
	case string:
		if v, err := ParseDateTime(time.DateOnly, val); err == nil {
			*d = NewDate(v)
		}
	case time.Time:
		if v, err := ParseDateTime(time.DateOnly, val.Format(time.DateOnly)); err == nil {
			*d = NewDate(v)
		}
	}

	return nil
}

func (d Date) Value() (driver.Value, error) {
	v := time.Time(d).In(LocalTimeZone())

	return v.Format(time.DateOnly), nil
}

// JSON
func (d Date) MarshalJSON() ([]byte, error) {
	v := time.Time(d).In(LocalTimeZone())

	return json.Marshal(v.Format(time.DateOnly))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	v, err := ParseDateTime(time.DateOnly, s)

	if err != nil {
		return err
	}

	*d = NewDate(v)

	return nil
}

// String
func (d Date) String() string {
	v := time.Time(d).In(LocalTimeZone())

	return v.Format(time.DateOnly)
}

// ---------------------------------------------------------
//
//  Time
//
// ---------------------------------------------------------

type Time time.Time

func NewTime(t *time.Time) Time {
	v := time.Now().In(LocalTimeZone())

	if t != nil {
		val := t.In(LocalTimeZone())

		return Time(
			time.Date(
				v.Year(), v.Month(), v.Day(), val.Hour(), val.Minute(), val.Second(), 0, v.Location(),
			),
		)
	} else {
		return Time(
			time.Date(
				v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second(), 0, v.Location(),
			),
		)
	}
}

// GORM
func (t Time) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	if db.Dialector.Name() == "sqlite" {
		return "TEXT"
	} else {
		return "TIME"
	}
}

func (t *Time) Scan(value any) error {
	switch val := value.(type) {
	case []byte:
		if v, err := ParseDateTime(time.TimeOnly, string(val)); err == nil {
			*t = NewTime(v)
		}
	case string:
		if v, err := ParseDateTime(time.TimeOnly, val); err == nil {
			*t = NewTime(v)
		}
	case time.Time:
		if v, err := ParseDateTime(time.TimeOnly, val.Format(time.TimeOnly)); err == nil {
			*t = NewTime(v)
		}
	}

	return nil
}

func (t Time) Value() (driver.Value, error) {
	v := time.Time(t).In(LocalTimeZone())

	return v.Format(time.TimeOnly), nil
}

// JSON
func (t Time) MarshalJSON() ([]byte, error) {
	v := time.Time(t).In(LocalTimeZone())

	return json.Marshal(v.Format(time.TimeOnly))
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	v, err := ParseDateTime(time.TimeOnly, s)

	if err != nil {
		return err
	}

	*t = NewTime(v)

	return nil
}

// String
func (t Time) String() string {
	v := time.Time(t).In(LocalTimeZone())

	return v.Format(time.TimeOnly)
}
