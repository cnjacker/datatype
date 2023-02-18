package pqarray

import (
	"database/sql/driver"

	"github.com/lib/pq"
)

type FloatArray []float64

// GORM
func (a FloatArray) GormDataType() string {
	return "float[]"
}

func (a *FloatArray) Scan(value any) error {
	var objs pq.Float64Array

	if err := objs.Scan(value); err == nil {
		*a = FloatArray(objs)
	}

	return nil
}

func (a FloatArray) Value() (driver.Value, error) {
	return pq.Float64Array(a).Value()
}

func (a FloatArray) Array() []float64 {
	return []float64(a)
}
