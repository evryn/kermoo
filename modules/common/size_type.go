package common

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/docker/go-units"
)

// Size represents bytes
type Size int64

func (s *Size) ToBytes() int64 {
	return int64(*s)
}

func (s Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Size) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*s = Size(value)
		return nil
	case string:
		tmp, err := units.FromHumanSize(value)
		if err != nil {
			return fmt.Errorf("size parsing failed: %v", err)
		}
		*s = Size(tmp)
		return nil
	default:
		return errors.New("invalid type for size")
	}
}
