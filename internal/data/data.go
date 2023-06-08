package data

import (
	"fmt"
	"strconv"
)

type DataString string

func (s *DataString) UnmarshalJSON(b []byte) error {
	value, err := strconv.Unquote(string(b))
	if err != nil {
		return fmt.Errorf("invalid json format")
	}
	*s = DataString(value)
	return nil
}

type Alert struct {
	Type    string
	Content string
}
