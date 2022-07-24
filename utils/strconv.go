package utils

import (
	"fmt"
	"strings"
)

// https://stackoverflow.com/questions/37532255/one-liner-to-transform-int-into-string
func UintArrayToString(a []uint64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
