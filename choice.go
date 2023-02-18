package datatype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Choice struct {
	Code  uint8  `json:"code" validate:"required"`
	Value string `json:"value"`
}

type ChoiceType struct {
	Choice Choice `json:"choice" validate:"required"`
	Meta   []Choice
}

// Update
func (c *ChoiceType) Update(code uint8) bool {
	ok := false

	if len(c.Meta) > 0 {
		for _, meta := range c.Meta {
			if meta.Code == code {
				c.Choice = meta
				ok = true

				break
			}
		}

		if !ok {
			c.Choice = c.Meta[0]
		}
	} else {
		c.Choice = Choice{Code: code}
		ok = true
	}

	return ok
}

// GORM
func (c *ChoiceType) Scan(value any) error {
	if v, ok := value.(int); ok {
		c.Update(uint8(v))
	}

	return nil
}

func (c ChoiceType) Value() (driver.Value, error) {
	return c.Choice.Code, nil
}

// JSON
func (c ChoiceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Choice)
}

func (c *ChoiceType) UnmarshalJSON(data []byte) error {
	var choice Choice

	if err := json.Unmarshal(data, &choice.Code); err != nil {
		if err := json.Unmarshal(data, &choice); err != nil {
			return err
		}
	}

	c.Update(choice.Code)

	return nil
}

// String
func (c ChoiceType) String() string {
	return fmt.Sprintf("{%d %s}", c.Choice.Code, c.Choice.Value)
}
