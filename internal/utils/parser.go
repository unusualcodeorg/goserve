package utils

import "strconv"

func ParseUint16(str string) uint16 {
	u, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return 0
	}
	return uint16(u)
}

func ParseUint8(str string) uint8 {
	u, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		return 0
	}
	return uint8(u)
}
