package types

import (
	"barlio/internal/helper"
	"database/sql/driver"
	"fmt"
)

type String string

func (s String) Value() (driver.Value, error) {
	if helper.StringIsNotEmpty(s) {
		return driver.Value(string(s)), nil
	}
	return driver.Value(nil), nil
}

func (s *String) Scan(i interface{}) error {
	switch i.(type) {
	case string:
		*s = String(i.(string))
		return nil
	case []uint8:
		*s = String(i.([]uint8))
		return nil
	case nil:
		*s = String("")
		return nil
	default:
		return fmt.Errorf("invalid type for string, %T", i)
	}
}

type Alert struct {
	Type    string
	Content string
}
