package config

import "testing"

func Test_Age_Lexical_Positive(t *testing.T) {
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

func Test_Age_Lexical_Negative(t *testing.T) {
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

// month-name, weekday-name, numeric-day, hour+"h", minute+"m", second+"s"
func Test_Range_Lexical(t *testing.T) {
	//years+"y" | months+"M" | weeks+"w" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
	validValues := [...]Range{"January", "Jan", "Jan.",
		" February", "Feb", "feb.", "FE", "fe",
		"March ", "Mar", "mar.", "MR", "mr",
		" April", "Apr", "apr.", "AP", " ap ",
		" May ", "MY", "my",
		" June ", "Jun", "jun.", "JN", " jn ",
		" July ", "Jul", "jul.", "JN", " jn ",
		" August ", "Aug", "aug. ", "AU", " au ",
		" September ", "Sep", "sep. ", "Sept", " sept. ", "SE", " se ",
		" October ", "Oct", "oct. ", "OC", " oc ",
		" November ", "Nov", "nov. ", "NV", " nv ",
		" December ", "Dec", "dec. ", "DE", " de ",
		//
		"Monday ", "Mon", "mon.", "mo", "MO.",
		"Tuesday ", "Tue", "tue.", "tu", "TU.",
		"Wednesday ", "Wed", "wed.", "we", "WE.",
		"Thursday ", "Thu", "thu.", "th", "TH.",
		"Friday  ", "Fri", "fri.", "fr", "FR.",
		"  Saturday  ", "Sat", "sat.", "sa", "SA.",
		"  Sunday", "Sun", "sun.", "su", "SU.",
		//
		" 31 dec 00h 00m 00s", " 31 dec 23h 59m 59s",
		" wed 10h", " wed 11h",
	}

	for _, v := range validValues {
		if err := v.HasError(); err != nil {
			t.Errorf("Expected valid lexical value for '%s': %s", v, err.Error())
		}
	}
}
