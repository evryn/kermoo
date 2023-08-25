package fluent

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func prepareUnmarshalString(data []byte) (string, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return "", err
	}

	var input string
	switch value := v.(type) {
	case string:
		input = value
	case float64:
		input = strconv.FormatFloat(value, 'f', -1, 64)
	case int:
		input = strconv.Itoa(value)
	default:
		return "", fmt.Errorf("unsupported fluent type: %T", v)
	}

	return input, nil
}
