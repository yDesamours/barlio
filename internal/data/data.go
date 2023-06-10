package data

import (
	"fmt"
	"strconv"
)

type String string

func (s *String) UnmarshalJSON(b []byte) error {
	value, err := strconv.Unquote(string(b))
	if err != nil {
		return fmt.Errorf("invalid json format")
	}
	*s = String(value)
	return nil
}

type Alert struct {
	Type    string
	Content string
}
