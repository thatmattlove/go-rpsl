package serialize

import (
	"fmt"
	"strings"
)

func StringsStringers[T ~[]E, E fmt.Stringer](stringers T) []string {
	out := make([]string, len(stringers))
	for i := range stringers {
		out[i] = stringers[i].String()
	}
	return out
}

func JoinStringers[T ~[]E, E fmt.Stringer](stringers T, sep string) string {
	return strings.Join(StringsStringers(stringers), sep)
}
