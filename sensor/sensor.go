// Package sensor reads temperature values reported by a sensor.
package sensor

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseCelsius parses one raw sensor line, such as "23.5", into degrees Celsius.
func ParseCelsius(raw string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return 0, fmt.Errorf("parse temperature %q: %w", raw, err)
	}
	return value, nil
}
