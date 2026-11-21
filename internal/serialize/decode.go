package serialize

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ProcessString(current string, _new []byte, sep string) string {
	parts := strings.Split(current, sep)
	newParts := strings.Split(string(_new), sep)
	parts = append(parts, newParts...)
	cleanParts := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			cleanParts = append(cleanParts, strings.TrimSpace(p))
		}
	}
	return strings.Trim((strings.Join(cleanParts, sep)), "\n")
}

func ProcessStringSlice(current reflect.Value, _new []byte) reflect.Value {
	for _, p := range strings.Split(string(_new), ",") {
		// Trim any surrounding whitespace on each member and skip empty members.
		if p = strings.TrimSpace(p); p != "" {
			current = reflect.Append(current, reflect.ValueOf(p))
		}
	}
	return current
}

// isUnmarshaler reports whether m has the signature func([]byte) (T, error), where T is
// assignable to t.
func isUnmarshaler(m, t reflect.Type) bool {
	return m.NumIn() == 1 && !m.IsVariadic() && m.In(0) == reflect.TypeFor[[]byte]() &&
		m.NumOut() == 2 && m.Out(0).AssignableTo(t) && m.Out(1) == reflect.TypeFor[error]()
}

// hasField reports whether t has a field with an rpsl tag matching key.
func hasField(t reflect.Type, key string) bool {
	for i := range t.NumField() {
		name, _, _ := strings.Cut(t.Field(i).Tag.Get("rpsl"), ",")
		if name != "-" && t.Field(i).IsExported() && strings.EqualFold(name, key) {
			return true
		}
	}
	return false
}

// Decode decodes a byte string of RPSL data to a Go RPSL object.
// The second argument must be a pointer to an RPSL struct.
func Decode(b []byte, o any) error {
	to := reflect.TypeOf(o)
	rv := reflect.ValueOf(o)
	// Ensure passed value is a non-nil pointer.
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &UnmarshalBinaryErr{to}
	}
	// Retrieve field pointer value.
	rt := to.Elem()
	// Ensure passed value is a struct.
	if rt.Kind() != reflect.Struct {
		return &UnmarshalBinaryErr{to}
	}
	// Retrieve field value's value.
	rvElem := rv.Elem()
	// Trim leading & trailing newlines.
	b = bytes.Trim(b, "\n")
	// Separate full blob by lines.
	blines := bytes.Split(b, []byte{0xa})
	// Remove the first line's indentation if all lines share it.
	indent := blines[0][:len(blines[0])-len(bytes.TrimLeft(blines[0], " \t"))]
	for _, line := range blines {
		if len(bytes.TrimSpace(line)) > 0 && !bytes.HasPrefix(line, indent) {
			indent = nil
		}
	}
	for i := range blines {
		blines[i] = bytes.TrimPrefix(blines[i], indent)
	}

	// Create a slice of key/value pairs.
	pairs := make([][][]byte, 0, len(blines))
	for i := range blines {
		if len(blines[i]) == 0 {
			continue
		}
		switch blines[i][0] {
		case '%', '#':
			// Skip comment lines.
			continue
		case ' ', '\t', '+':
			// Append continuation line to the previous pair's values, joined with a single space.
			cont := bytes.TrimSpace(blines[i][1:])
			parsed, _, _ := bytes.Cut(cont, []byte{0x23})
			if len(pairs) > 0 {
				prev := pairs[len(pairs)-1]
				prev[1] = bytes.TrimSpace(bytes.Join([][]byte{prev[1], cont}, []byte{0x20}))
				prev[2] = bytes.TrimSpace(bytes.Join([][]byte{prev[2], parsed}, []byte{0x20}))
			}
			continue
		}
		// Split each line on the first ':' (accounts for ':' in value).
		k, v, found := bytes.Cut(blines[i], []byte{0x3a})
		if !found {
			// If no ':' character found, skip this line
			continue
		}
		// Trim any surrounding whitespace on key.
		key := bytes.TrimSpace(k)
		// Trim any surrounding whitespace on value.
		value := bytes.TrimSpace(v)
		// Strip any trailing '#' comment from values that are parsed.
		parsed, _, _ := bytes.Cut(value, []byte{0x23})
		// Add pair to k/v pair slice.
		pairs = append(pairs, [][]byte{key, value, bytes.TrimSpace(parsed)})
	}
	for i := range rt.NumField() {
		// Retrieve struct field.
		field := rt.Field(i)
		// Retrieve struct field's value.
		valueField := rvElem.Field(i)

		tag := field.Tag.Get("rpsl")
		if tag == "" || !field.IsExported() {
			// If no rpsl tag is present or the field is unexported, skip this field (hidden or
			// irrelevant field).
			continue
		}
		// Get rpsl struct tag value, ignoring any ',omitempty' tags.
		tags := strings.Split(tag, ",")
		keyName := tags[0]

		as, hasAs := field.Tag.Lookup("as")

		// Add extra values to the 'Extra' field map, which is tagged as "-".
		if keyName == "-" {
			m, ok := valueField.Interface().(map[string]string)
			if !ok {
				return fmt.Errorf("field '%s' has an invalid rpsl tag", field.Name)
			}
			if m == nil {
				// Initialize map so it is never nil after decoding.
				m = make(map[string]string)
				valueField.Set(reflect.ValueOf(m))
			}
			for _, pair := range pairs {
				// Skip attributes that have a matching struct field.
				if hasField(rt, string(pair[0])) {
					continue
				}
				// Join repeated attributes by newline.
				m[string(pair[0])] = ProcessString(m[string(pair[0])], pair[1], "\n")
			}
			continue
		}

		// Begin struct field to key/pair matching.
		for _, pair := range pairs {
			key := pair[0]    // left side of first ':', key
			value := pair[1]  // right side of first ':', value
			parsed := pair[2] // value without '#' comments

			// Struct field matches this k/v pair. Attribute names are case-insensitive.
			if strings.EqualFold(keyName, string(key)) {
				if hasAs {
					switch valueType := valueField.Interface().(type) {
					case string:
						switch as {
						case "multiline":
							valueField.SetString(ProcessString(valueType, value, "\n"))
						case "comma-space":
							valueField.SetString(ProcessString(valueType, value, ", "))
						case "comma":
							valueField.SetString(ProcessString(valueType, value, ","))
						}
					case []string:
						// List values may be comma-separated regardless of how they are encoded.
						valueField.Set(ProcessStringSlice(valueField, parsed))
					}
				} else {
					// If the value is a type that has a valid UnmarshalBinary method, call that method
					// and assign the value the field.
					if method := valueField.MethodByName("UnmarshalBinary"); method.IsValid() && isUnmarshaler(method.Type(), field.Type) {
						args := []reflect.Value{reflect.ValueOf(parsed)}
						result := method.Call(args)
						// Check and coerce to error if error is non-nil.
						maybeErr := result[1].Interface()
						if maybeErr != nil {
							err := maybeErr.(error)
							err = errors.Join(fmt.Errorf("%s: failed to unmarshal", keyName), err)
							return err
						}
						// Assign return value from UnmarshalBinary to struct field.
						valueField.Set(result[0])
						continue
					}
					// If the value is a (supported) native type, handle accordingly:
					switch valueField.Kind() {
					case reflect.String:
						valueField.SetString(string(parsed))
						continue
					case reflect.Uint32:
						u, err := strconv.ParseUint(string(parsed), 10, 32)
						if err != nil {
							err := errors.Join(fmt.Errorf("%s: value could not be parsed as uint32", keyName), err)
							return err
						}
						valueField.SetUint(u)
						continue
					}
				}
			}
		}
	}
	return nil
}
