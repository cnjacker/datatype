package datatype

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// ---------------------------------------------------------
//
//  JSON
//
// ---------------------------------------------------------

type JSON map[string]any

// GORM
func (j JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	if db.Dialector.Name() == "sqlserver" {
		return "NVARCHAR(MAX)"
	} else {
		return "JSON"
	}
}

func (j JSON) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(j) > 0 {
		if v, err := json.Marshal(j); err == nil {
			return gorm.Expr("?", string(v))
		}
	}

	return gorm.Expr("NULL")
}

func (j *JSON) Scan(value any) error {
	if value != nil {
		var bytes []byte

		switch v := value.(type) {
		case []byte:
			bytes = v
		case string:
			bytes = []byte(v)
		}

		if len(bytes) > 0 {
			return json.Unmarshal(bytes, j)
		}
	}

	*j = nil

	return nil
}

// String
func (j JSON) String() string {
	if v, err := json.MarshalIndent(j, "", "    "); err == nil {
		return string(v)
	}

	return ""
}

// ---------------------------------------------------------
//
//  JSONQuery
//
// ---------------------------------------------------------

type JSONQueryExpression struct {
	column      string
	keys        []string
	hasKeys     bool
	equals      bool
	equalsValue any
	extract     bool
	path        string
}

// Query
func JSONQuery(column string) *JSONQueryExpression {
	return &JSONQueryExpression{column: column}
}

// Extract
func (jsonQuery *JSONQueryExpression) Extract(path string) *JSONQueryExpression {
	jsonQuery.extract = true
	jsonQuery.path = path

	return jsonQuery
}

// HasKey
func (jsonQuery *JSONQueryExpression) HasKey(keys ...string) *JSONQueryExpression {
	jsonQuery.keys = keys
	jsonQuery.hasKeys = true

	return jsonQuery
}

// Equals
func (jsonQuery *JSONQueryExpression) Equals(value any, keys ...string) *JSONQueryExpression {
	jsonQuery.keys = keys
	jsonQuery.equals = true
	jsonQuery.equalsValue = value

	return jsonQuery
}

// Build
func (jsonQuery *JSONQueryExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		switch stmt.Dialector.Name() {
		case "mysql", "sqlite":
			switch {
			case jsonQuery.extract:
				builder.WriteString("JSON_EXTRACT(")
				builder.WriteQuoted(jsonQuery.column)
				builder.WriteByte(',')
				builder.AddVar(stmt, jsonQuery.path)
				builder.WriteString(")")
			case jsonQuery.hasKeys:
				if len(jsonQuery.keys) > 0 {
					builder.WriteString("JSON_EXTRACT(")
					builder.WriteQuoted(jsonQuery.column)
					builder.WriteByte(',')
					builder.AddVar(stmt, jsonQueryJoin(jsonQuery.keys))
					builder.WriteString(") IS NOT NULL")
				}
			case jsonQuery.equals:
				if len(jsonQuery.keys) > 0 {
					builder.WriteString("JSON_EXTRACT(")
					builder.WriteQuoted(jsonQuery.column)
					builder.WriteByte(',')
					builder.AddVar(stmt, jsonQueryJoin(jsonQuery.keys))
					builder.WriteString(") = ")
					if value, ok := jsonQuery.equalsValue.(bool); ok {
						builder.WriteString(strconv.FormatBool(value))
					} else {
						stmt.AddVar(builder, jsonQuery.equalsValue)
					}
				}
			}
		case "postgres":
			switch {
			case jsonQuery.extract:
				builder.WriteString(fmt.Sprintf("json_extract_path_text(%v::json,", stmt.Quote(jsonQuery.column)))
				stmt.AddVar(builder, jsonQuery.path)
				builder.WriteByte(')')
			case jsonQuery.hasKeys:
				if len(jsonQuery.keys) > 0 {
					stmt.WriteQuoted(jsonQuery.column)
					stmt.WriteString("::json")
					for _, key := range jsonQuery.keys[0 : len(jsonQuery.keys)-1] {
						stmt.WriteString(" -> ")
						stmt.AddVar(builder, key)
					}

					stmt.WriteString(" ? ")
					stmt.AddVar(builder, jsonQuery.keys[len(jsonQuery.keys)-1])
				}
			case jsonQuery.equals:
				if len(jsonQuery.keys) > 0 {
					builder.WriteString(fmt.Sprintf("json_extract_path_text(%v::json,", stmt.Quote(jsonQuery.column)))

					for idx, key := range jsonQuery.keys {
						if idx > 0 {
							builder.WriteByte(',')
						}
						stmt.AddVar(builder, key)
					}
					builder.WriteString(") = ")

					if _, ok := jsonQuery.equalsValue.(string); ok {
						stmt.AddVar(builder, jsonQuery.equalsValue)
					} else {
						stmt.AddVar(builder, fmt.Sprint(jsonQuery.equalsValue))
					}
				}
			}
		}

	}
}

// jsonQueryJoin
func jsonQueryJoin(keys []string) string {
	if len(keys) == 1 {
		return "$." + keys[0]
	}

	n := len("$.")
	n += len(keys) - 1

	for i := 0; i < len(keys); i++ {
		n += len(keys[i])
	}

	var builder strings.Builder

	builder.Grow(n)
	builder.WriteString("$.")
	builder.WriteString(keys[0])

	for _, key := range keys[1:] {
		builder.WriteString(".")
		builder.WriteString(key)
	}

	return builder.String()
}
