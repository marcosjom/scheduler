package config

import "testing"

// Testing some valid examples
func Test_Age_Lexical_Valid(t *testing.T) {
	validValues := [...]Age{"1y 2M 3h 22m 99s",
		"46977y", "20y",
		"11M", "88M",
		"40h", "5h",
		"1069m", "29m",
		"1789s", "60s", "1s",
	}
	for _, v := range validValues {
		if err := v.HasError(); err != nil {
			t.Errorf("Expected valid lexical value for '%s': %s", v, err.Error())
		}
	}
}

// Testing some invalid examples
func Test_Age_Lexical_Invalid(t *testing.T) {
	validValues := [...]Age{"", "0", "0.0", "20", "-1", "1.1", "0.1", "-1.1", "-0.1",
		"1Y 2mm 3hh 22mm 99ss",
		"-1y", "0y",
		"-11M", "88MM",
		"-4.0h", "-5.h",
		"1,069.m", "2.9m",
		"1,789.s", "6.0s", "1.s",
	}
	for _, v := range validValues {
		if err := v.HasError(); err == nil {
			t.Errorf("Expected invalid lexical value for '%s'.", v)
		}
	}
}
