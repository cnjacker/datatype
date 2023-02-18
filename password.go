package datatype

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hash []byte
}

// Update
func (p *Password) Update(password string) {
	p.hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Compare
func (p *Password) Compare(password string) bool {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(password)) == nil
}

// GORM
func (p *Password) Scan(value any) error {
	if v, ok := value.([]byte); ok {
		p.hash = v
	}

	return nil
}

func (p Password) Value() (driver.Value, error) {
	return p.hash, nil
}

// JSON
func (p Password) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.Repeat("*", 8))
}

func (p *Password) UnmarshalJSON(data []byte) error {
	return nil
}

// String
func (p Password) String() string {
	return string(p.hash)
}
