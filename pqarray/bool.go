package pqarray

import (
	"database/sql/driver"

	"github.com/lib/pq"
)

type BoolArray []bool

// GORM
func (a BoolArray) GormDataType() string {
	return "bool[]"
}

func (a *BoolArray) Scan(value any) error {
	var objs pq.BoolArray

	if err := objs.Scan(value); err == nil {
		*a = BoolArray(objs)
	}

	return nil
}

func (a BoolArray) Value() (driver.Value, error) {
	return pq.BoolArray(a).Value()
}

func (a BoolArray) Array() []bool {
	return []bool(a)
}
