package config

import "testing"

func Test_Command_Lexical_Positive(t *testing.T) {
	validValues := [...]Command{Command{Execute: "do somthing"},
		Command{Execute: "exec/bin/program with params"},
	}
	for _, v := range validValues {
		if err := v.HasError(); err != nil {
			t.Errorf("Expected valid lexical value for '%s': %s", v, err.Error())
		}
	}
}

func Test_Command_Lexical_Negative(t *testing.T) {
	validValues := [...]Command{Command{Execute: ""},
		Command{Execute: "   "},
		Command{Execute: "\t\r\n"},
	}
	for _, v := range validValues {
		if err := v.HasError(); err == nil {
			t.Errorf("Expected invalid lexical value for '%s'", v)
		}
	}
}
