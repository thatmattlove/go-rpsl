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

// UnmarshalBinary parses a byte string to a description, maintaining the previous descr
// lines when applicable.
func (d Description) UnmarshalBinary(b []byte) (Description, error) {
	asS := string(b)
	if len(d) == 0 {
		return Description(asS), nil
	}
	parts := strings.Split(string(d), "\n")
	parts = append(parts, asS)
	return Description(strings.Join(parts, "\n")), nil
}
