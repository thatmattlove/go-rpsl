package rpsl

import (
	"fmt"
	"strings"
)

// The descr field can be a multiline value where each line is represented as a separate 'descr'
// key/value pair.
type Description string

// RPSL serializes the multiline description field.
func (d Description) RPSL() string {
	parts := strings.Split(string(d), "\n")
	out := ""
	for i := range parts {
		out += fmt.Sprintf("descr: %s\n", parts[i])
	}
	return strings.TrimSuffix(out, "\n")
}
