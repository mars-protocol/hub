package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// UintArrayToString joins an array of uint64 numbers into a string with the given delimiter
//
// Forked from https://stackoverflow.com/questions/37532255/one-liner-to-transform-int-into-string
func UintArrayToString(uints []uint64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(uints), " ", delim, -1), "[]")
}

// StringToUintArray parses a string with the given delimiter into an array of uint64 numbers
func StringToUintArray(str, delim string) ([]uint64, error) {
	uintStrs := strings.Split(str, delim)
	uints := []uint64{}
	for _, idStr := range uintStrs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ids: %s", err)
		}

		uints = append(uints, id)
	}

	return uints, nil
}
