package rpsl

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// UnmarshalBinary marshals a byte string of RPSL data to a Go RPSL object.
// The second argument must be a pointer to an RPSL struct.
//
// Example:
//
//	b := []byte(`aut-num: 65000`)
//	var autNum rpsl.AutNum
//	err := rpsl.UnmarshalBinary(b, &autNum)
func UnmarshalBinary(b []byte, o any) error {
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

	// Create a slice of key/value pairs.
	pairs := make([][][]byte, 0, len(blines))
	for i := range blines {
		// Split each line by ':'
		pairline := bytes.Split(blines[i], []byte{0x3a})
		if len(pairline) < 2 {
			// If no ':' character found, skip this line
			continue
		}
		// Trim any surrounding whitespace on key.
		key := bytes.TrimSpace(pairline[0])
		// Trim any surrounding whitespace on value.
		value := bytes.TrimSpace(
			// Join remaining characters (accounts for ':' in value)
			bytes.Join(pairline[1:], []byte{0x0}),
		)
		// Add pair to k/v pair slice.
		pairs = append(pairs, [][]byte{key, value})
	}
	for i := range rt.NumField() {
		// Retrieve struct field.
		field := rt.Field(i)
		// Retrieve struct field's value.
		valueField := rvElem.Field(i)

		tag := field.Tag.Get("rpsl")
		if tag == "" {
			// If no rpsl tag is present, skip this field (hidden or irrelevant field).
			continue
		}
		// Get rpsl struct tag value, ignoring any ',omitempty' tags.
		tags := strings.Split(tag, ",")
		keyName := tags[0]

		// Begin struct field to key/pair matching.
		for _, pair := range pairs {
			key := pair[0]   // left side of first ':', key
			value := pair[1] // right side of first ':', value
			if keyName == "-" {
				// Add extra values to the 'Extra' field map, which is tagged as "-".
				m := valueField.Interface().(map[string]string)
				if m == nil {
					// Initialize map if this is the first pass.
					m = make(map[string]string)
				}
				m[string(key)] = string(value)
				// Reassign 'Extra' map.
				valueField.Set(reflect.ValueOf(m))
				continue
			}
			// Struct field matches this k/v pair.
			if keyName == string(key) {
				// If the value is a type that has a valid UnmarshalBinary method, call that method
				// and assign the value the field.
				if method := valueField.MethodByName("UnmarshalBinary"); method.IsValid() {
					args := []reflect.Value{reflect.ValueOf(value)}
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
					valueField.SetString(string(value))
					continue
				case reflect.Uint32:
					u, err := strconv.ParseUint(string(value), 10, 64)
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
	return nil
}
