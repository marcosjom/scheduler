package common

import "strconv"

type Minute uint8

const Minute_invalid = Minute(0)

func (d Minute) String() string {
	if d >= 1 && d <= 61 {
		return strconv.Itoa(int(d) - 1)
	}
	return ""
}

func (d Minute) Index() uint8 {
	if d >= 1 && d <= 61 {
		return uint8(d) - 1
	}
	return uint8(Minute_invalid)
}
