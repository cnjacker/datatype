package datatype

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Coordinate struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
}

// GORM
func (c Coordinate) GormDataType() string {
	return "geometry"
}

func (c Coordinate) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []any{fmt.Sprintf("SRID=4326;POINT(%f %f)", c.Longitude, c.Latitude)},
	}
}

func (c *Coordinate) Scan(value any) error {
	if v, ok := value.(string); ok {
		if bytes, err := hex.DecodeString(v); err == nil {
			if geom, err := ewkb.Unmarshal(bytes); err == nil {
				if len(geom.FlatCoords()) >= 2 {
					c.Longitude = geom.FlatCoords()[0]
					c.Latitude = geom.FlatCoords()[1]
				}
			}
		}
	}

	return nil
}
