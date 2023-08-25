package fluent

import (
	"errors"
	"strconv"
	"strings"
	"time"
	// Assuming "units" package is available
)

type Parser struct {
	input string
}

// GetFloats parses the input string to extract an array of floating point numbers.
// The function recognizes single values, ranges ("min to max"), and arrays separated by commas.
func (p *Parser) GetFloats() (*ParsedValue[float64], error) {
	// Parsing ranges
	if strings.Contains(p.input, " to ") {
		parts := strings.Split(p.input, " to ")
		if len(parts) != 2 {
			return nil, errors.New("invalid range format. it must be in the form of \"min to max\" like \"1.5 to 6\"")
		}

		start, err := p.convertFloat(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}

		end, err := p.convertFloat(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, err
		}

		pv := newParsedValue[float64]([]float64{start, end}, true)

		return &pv, nil
	}

	// Parsing arrays
	if strings.Contains(p.input, ",") {
		var values []float64
		parts := strings.Split(p.input, ",")
		for _, part := range parts {
			val, err := p.convertFloat(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}

			values = append(values, val)
		}

		pv := newParsedValue[float64](values, false)

		return &pv, nil
	}

	// Parsing exact values
	val, err := p.convertFloat(strings.TrimSpace(p.input))
	if err != nil {
		return nil, err
	}

	pv := newParsedValue[float64]([]float64{val}, false)

	return &pv, nil
}

// GetSizes parses the input string to extract an array of sizes.
// The sizes can be in bytes, KiB, MiB, etc. The function recognizes single values, ranges, and comma-separated arrays.
func (p *Parser) GetSizes() (*ParsedValue[int64], error) {
	// Parsing ranges
	if strings.Contains(p.input, " to ") {
		parts := strings.Split(p.input, " to ")
		if len(parts) != 2 {
			return nil, errors.New("invalid range format. it must be in the form of \"min to max\" like \"1.5 to 6\"")
		}

		start, err := p.convertSize(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}

		end, err := p.convertSize(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, err
		}

		pv := newParsedValue[int64]([]int64{start, end}, true)

		return &pv, nil
	}

	// Parsing arrays
	if strings.Contains(p.input, ",") {
		var values []int64
		parts := strings.Split(p.input, ",")
		for _, part := range parts {
			val, err := p.convertSize(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}

			values = append(values, val)
		}

		pv := newParsedValue[int64](values, false)

		return &pv, nil
	}

	// Parsing exact values
	val, err := p.convertSize(strings.TrimSpace(p.input))
	if err != nil {
		return nil, err
	}

	pv := newParsedValue[int64]([]int64{val}, false)

	return &pv, nil
}

// GetDurations parses the input string to extract an array of durations.
// Recognizes single values, ranges, and comma-separated arrays.
func (p *Parser) GetDuations() (*ParsedValue[time.Duration], error) {
	// Parsing ranges
	if strings.Contains(p.input, " to ") {
		parts := strings.Split(p.input, " to ")
		if len(parts) != 2 {
			return nil, errors.New("invalid range format. it must be in the form of \"min to max\" like \"1.5 to 6\"")
		}

		start, err := p.convertDuration(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}

		end, err := p.convertDuration(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, err
		}

		pv := newParsedValue[time.Duration]([]time.Duration{start, end}, true)

		return &pv, nil
	}

	// Parsing arrays
	if strings.Contains(p.input, ",") {
		var values []time.Duration
		parts := strings.Split(p.input, ",")
		for _, part := range parts {
			val, err := p.convertDuration(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}

			values = append(values, val)
		}

		pv := newParsedValue[time.Duration](values, false)

		return &pv, nil
	}

	// Parsing exact values
	val, err := p.convertDuration(strings.TrimSpace(p.input))
	if err != nil {
		return nil, err
	}

	pv := newParsedValue[time.Duration]([]time.Duration{val}, false)

	return &pv, nil
}

// convertDuration converts a string representation of duration into a time.Duration type.
func (p *Parser) convertDuration(part string) (time.Duration, error) {
	return time.ParseDuration(part)
}

// convertFloat converts a string representation into a floating point number.
func (p *Parser) convertFloat(part string) (float64, error) {
	return strconv.ParseFloat(part, 64)
}

// convertSize converts a string representation of size into its equivalent in bytes.
// The function accommodates for various size suffixes like K, M, G, etc.
func (p *Parser) convertSize(size string) (int64, error) {
	if value, err := strconv.ParseInt(size, 10, 64); err == nil {
		return value, nil
	}

	var multiplier int64
	size = strings.TrimSpace(size)

	if strings.HasSuffix(size, "Ki") {
		multiplier = 1024
		size = strings.TrimSuffix(size, "Ki")
	} else if strings.HasSuffix(size, "K") {
		multiplier = 1000
		size = strings.TrimSuffix(size, "K")
	} else if strings.HasSuffix(size, "Mi") {
		multiplier = 1024 * 1024
		size = strings.TrimSuffix(size, "Mi")
	} else if strings.HasSuffix(size, "M") {
		multiplier = 1000 * 1000
		size = strings.TrimSuffix(size, "M")
	} else if strings.HasSuffix(size, "Gi") {
		multiplier = 1024 * 1024 * 1024
		size = strings.TrimSuffix(size, "Gi")
	} else if strings.HasSuffix(size, "G") {
		multiplier = 1000 * 1000 * 1000
		size = strings.TrimSuffix(size, "G")
	} else if strings.HasSuffix(size, "Ti") {
		multiplier = 1024 * 1024 * 1024 * 1024
		size = strings.TrimSuffix(size, "Ti")
	} else if strings.HasSuffix(size, "T") {
		multiplier = 1000 * 1000 * 1000 * 1000
		size = strings.TrimSuffix(size, "T")
	} else {
		return 0, errors.New("invalid syntax for size")
	}

	value, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		return 0, err
	}

	return value * multiplier, nil
}

// newParser initializes and returns a Parser object with the provided input string.
func newParser(input string) Parser {
	return Parser{
		input: input,
	}
}
