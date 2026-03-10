package config

import (
	"errors"
	"strconv"
	"strings"
)

// Age represents an ammount of time.
// years+"y" | months+"M" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
type Age string

// Lexical validation.
func (t Age) HasError() error {
	partsCount := 0
	parts := strings.Split(string(t), " ")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// Eval years+"y" | months+"M" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
		if len(p) == 1 {
			return errors.New("Part requires value and suffix: '" + p + "'.")
		}
		sfx := p[len(p)-1]
		switch sfx {
		case 'y', 'M', 'd', 'h', 'm', 's':
			num, err := strconv.Atoi(p[:len(p)-1])
			if err != nil {
				return errors.New("Part requires a numeric value: '" + p + "'.")
			} else if num <= 0 {
				return errors.New("Part requires a positive non-zero value: '" + p + "'.")
			}
			partsCount++
		default:
			return errors.New("Invalid suffix '" + string(sfx) + "'.")
		}
	}
	//
	if partsCount <= 0 {
		return errors.New("Empty age's value.")
	}
	//
	return nil
}
