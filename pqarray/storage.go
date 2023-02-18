package pqarray

import (
	"database/sql/driver"

	"github.com/cnjacker/datatype"
	"github.com/lib/pq"
)

type StorageArray []datatype.Storage

// GORM
func (a StorageArray) GormDataType() string {
	return "varchar[]"
}

func (a *StorageArray) Scan(value any) error {
	var objs pq.StringArray

	if err := objs.Scan(value); err == nil {
		var v []datatype.Storage

		for _, obj := range objs {
			v = append(v, datatype.Storage(datatype.Storage(obj).BindSignature()))
		}

		*a = v
	}

	return nil
}

func (a StorageArray) Value() (driver.Value, error) {
	var objs pq.StringArray

	for _, obj := range a {
		objs = append(objs, obj.UnBindSignature())
	}

	return objs.Value()
}

func (a StorageArray) Array() []string {
	var v []string

	for _, obj := range a {
		v = append(v, obj.UnBindSignature())
	}

	return v
}
