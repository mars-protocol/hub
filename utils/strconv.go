package utils

import (
	"fmt"
	"strings"
)

// UintArrayToString joins an array of uint64 numbers into a string with the
// given delimiter.
//
// Forked from https://stackoverflow.com/questions/37532255/one-liner-to-transform-int-into-string
func UintArrayToString(uints []uint64, delim string) string {
	return strings.Trim(strings.ReplaceAll(fmt.Sprint(uints), " ", delim), "[]")
}
