package pqarray

import (
	"database/sql/driver"

	"github.com/lib/pq"
)

type StringArray []string

// GORM
func (a StringArray) GormDataType() string {
	return "varchar[]"
}

func (a *StringArray) Scan(value any) error {
	var objs pq.StringArray

	if err := objs.Scan(value); err == nil {
		*a = StringArray(objs)
	}

	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	return pq.StringArray(a).Value()
}

func (a StringArray) Array() []string {
	return []string(a)
}
