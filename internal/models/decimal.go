package models

import (
	"database/sql/driver"

	"github.com/ericlagergren/decimal"
)

type Decimal struct {
	decimal.Big
}

func NewDecimal(value int64, scale int) *Decimal {
	return &Decimal{
		Big: *decimal.New(value, scale),
	}
}

func (d Decimal) Value() (driver.Value, error) {
	return d.Big.Quantize(6).String(), nil
}

func (d *Decimal) Scan(src interface{}) error {
	switch v := src.(type) {
	case int64:
		d = NewDecimal(v, 0)
	case []byte:

		return d.UnmarshalText(v)
	case string:

		return d.UnmarshalText([]byte(v))
	}

	return nil
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(d.Big.Quantize(6).String()), nil
}
