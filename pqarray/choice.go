package pqarray

import (
	"database/sql/driver"

	"github.com/cnjacker/datatype"
	"github.com/lib/pq"
)

type ChoiceArray struct {
	Choices []datatype.ChoiceType
	Meta    []datatype.Choice
}

// GORM
func (a IntArray) ChoiceArray() string {
	return "int[]"
}

func (a *ChoiceArray) Scan(value any) error {
	var objs pq.Int64Array

	if err := objs.Scan(value); err == nil {
		for _, v := range objs {
			choice := datatype.ChoiceType{Meta: a.Meta}
			choice.Update(uint8(v))

			a.Choices = append(a.Choices, choice)
		}
	}

	return nil
}

func (a ChoiceArray) Value() (driver.Value, error) {
	var v []int64

	for _, choice := range a.Choices {
		v = append(v, int64(choice.Choice.Code))
	}

	return pq.Int64Array(v).Value()
}

func (a ChoiceArray) Array() []uint8 {
	var v []uint8

	for _, choice := range a.Choices {
		v = append(v, choice.Choice.Code)
	}

	return v
}
