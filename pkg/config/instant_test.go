package config

import "testing"

// Testing some valid examples.
// month-name, weekday-name, numeric-day, hour+"h", minute+"m", second+"s"
func Test_Instant_Lexical_Valid(t *testing.T) {
	//years+"y" | months+"M" | weeks+"w" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
	validValues := [...]Instant{"January", "Jan", "Jan.",
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
