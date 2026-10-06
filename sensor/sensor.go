// Package sensor reads temperature values reported by a sensor.
package sensor

import (
	"fmt"
	"strconv"
	"strings"
)

// MaxCelsius is the temperature threshold above which an over-temperature condition exists.
const MaxCelsius = 80.0

// ParseCelsius parses one raw sensor line, such as "23.5", into degrees Celsius.
func ParseCelsius(raw string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return 0, fmt.Errorf("parse temperature %q: %w", raw, err)
	}
	return value, nil
}

// IsOverTemp reports whether the given temperature in degrees Celsius exceeds the safe threshold.
func IsOverTemp(celsius float64) bool {
	return celsius > MaxCelsius
}
