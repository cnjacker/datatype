package pqarray

import (
	"database/sql/driver"

	"github.com/lib/pq"
)

type IntArray []int64

// GORM
func (a IntArray) GormDataType() string {
	return "int[]"
}

func (a *IntArray) Scan(value any) error {
	var objs pq.Int64Array

	if err := objs.Scan(value); err == nil {
		*a = IntArray(objs)
	}

	return nil
}

func (a IntArray) Value() (driver.Value, error) {
	return pq.Int64Array(a).Value()
}

func (a IntArray) Array() []int64 {
	return []int64(a)
}
