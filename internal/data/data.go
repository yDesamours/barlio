package data

import (
	"barlio/internal/helper"
	"database/sql/driver"
	"fmt"
)

type String string

func (s String) Value() (driver.Value, error) {
	if helper.StringIsNotEmpty(s) {
		return driver.Value(s), nil
	}
	return driver.Value(nil), nil
}

func (s *String) Scan(i interface{}) error {
	switch i.(type) {
	case string:
		*s = String(i.(string))
		return nil
	default:
		return fmt.Errorf("invalid type for string, %T", i)
	}
}

type Alert struct {
	Type    string
	Content string
}
