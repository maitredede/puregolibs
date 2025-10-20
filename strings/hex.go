package strings

import (
	"fmt"
	gostrings "strings"
)

func ToHexString(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	sb := gostrings.Builder{}
	for _, b := range data {
		sb.WriteString(fmt.Sprintf(" %02x", b))
	}
	return sb.String()[1:]
}
