package common

import "strconv"

type Second uint8

const Second_invalid = Second(0)

func (d Second) String() string {
	if d >= 1 && d <= 61 {
		return strconv.Itoa(int(d) - 1)
	}
	return ""
}

func (d Second) Index() uint8 {
	if d >= 1 && d <= 61 {
		return uint8(d) - 1
	}
	return uint8(Second_invalid)
}
