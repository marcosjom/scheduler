package common

import "strconv"

type Monthday uint8

const Monthday_invalid = Monthday(0)

func (d Monthday) String() string {
	if d >= 1 && d <= 31 {
		return strconv.Itoa(int(d))
	}
	return ""
}

func (d Monthday) Index() uint8 {
	if d >= 1 && d <= 31 {
		return uint8(d) - 1
	}
	return uint8(Monthday_invalid)
}
