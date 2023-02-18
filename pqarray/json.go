package pqarray

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/cnjacker/datatype"
	"github.com/lib/pq"
)

type JSONArray []datatype.JSON

// GORM
func (a JSONArray) GormDataType() string {
	return "json[]"
}

func (a *JSONArray) Scan(value any) error {
	var objs pq.StringArray

	if err := objs.Scan(value); err == nil {
		var m []datatype.JSON

		for _, obj := range objs {
			var v datatype.JSON

			if err := json.Unmarshal([]byte(obj), &v); err == nil {
				m = append(m, v)
			}
		}

		*a = JSONArray(m)
	}

	return nil
}

func (a JSONArray) Value() (driver.Value, error) {
	var objs pq.StringArray

	for _, obj := range a {
		v, err := json.Marshal(obj)

		if err != nil {
			return nil, err
		}

		objs = append(objs, string(v))
	}

	return objs.Value()
}

func (a JSONArray) Array() []datatype.JSON {
	return []datatype.JSON(a)
}
