package pqarray

import (
	"database/sql/driver"

	"github.com/lib/pq"
)

type ByteArray [][]byte

// GORM
func (a ByteArray) GormDataType() string {
	return "bytea[]"
}

func (a *ByteArray) Scan(value any) error {
	var objs pq.ByteaArray

	if err := objs.Scan(value); err == nil {
		*a = ByteArray(objs)
	}

	return nil
}

func (a ByteArray) Value() (driver.Value, error) {
	return pq.ByteaArray(a).Value()
}

func (a ByteArray) Array() [][]byte {
	return [][]byte(a)
}
