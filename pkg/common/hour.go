package common

import "strconv"

type Hour uint8

const Hour_invalid = Hour(0)

func (d Hour) String() string {
	if d >= 1 && d <= 25 {
		return strconv.Itoa(int(d) - 1)
	}
	return ""
}

func (d Hour) Index() uint8 {
	if d >= 1 && d <= 25 {
		return uint8(d) - 1
	}
	return uint8(Hour_invalid)
}
