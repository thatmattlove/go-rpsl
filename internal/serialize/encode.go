package serialize

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

var ErrMustBeStruct = errors.New("value must be a struct")

func ProcessAsString(key, val string) string {
	return fmt.Sprintf("%s: %s\n", key, strings.TrimSpace(strings.Trim(val, "\n")))
}

func ProcessAsStringSlice(key string, val []string, sep string) string {
	vs := strings.TrimSpace(strings.ReplaceAll(strings.Join(val, sep), "\n", " "))
	return fmt.Sprintf("%s: %s\n", key, vs)
}

func ProcessAsMultilineString(key, val string) string {
	out := ""
	for _, p := range strings.Split(val, "\n") {
		if p != "" {
			out += ProcessAsString(key, p)
		}
	}
	return out
}

func ProcessAsMultilineStringSlice(key string, val []string) string {
	out := ""
	for _, p := range val {
		if p != "" {
			out += ProcessAsString(key, p)
		}
	}
	return out
}

// Encode encodes an RPSL object to a byte string. The argument must be a pointer to an RPSL
// object.
func Encode(s any) ([]byte, error) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrMustBeStruct
	}
	out := ""
	for i := range t.NumField() {
		field := t.Field(i)
		valueField := v.Field(i)
		sval := valueField.Interface()
		tag := field.Tag.Get("rpsl")
		if tag == "" {
			continue
		}
		if tag == "-" {
			for k, v := range sval.(map[string]string) {
				if k != "" && v != "" {
					out += fmt.Sprintf("%s: %s\n", k, v)
				}
			}
			continue
		}
		tags := strings.Split(tag, ",")
		omitEmpty := slices.Contains(tags, "omitempty")
		if omitEmpty && valueField.IsZero() {
			continue
		}
		key := tags[0]
		if key == "" {
			return nil, fmt.Errorf("field '%s' has an invalid rpsl tag", field.Name)
		}
		as, hasAs := field.Tag.Lookup("as")
		if hasAs {
			switch stype := sval.(type) {
			case string:
				switch as {
				case "multiline":
					out += ProcessAsMultilineString(key, stype)
				default:
					out += ProcessAsString(key, stype)
				}
			case []string:
				switch as {
				case "multiline":
					out += ProcessAsMultilineStringSlice(key, stype)
				case "comma-space":
					out += ProcessAsStringSlice(key, stype, ", ")
				case "comma":
					out += ProcessAsStringSlice(key, stype, ",")
				}
			}
		} else {
			value := ""
			switch stype := sval.(type) {
			case uint32:
				value = strconv.FormatUint(uint64(stype), 10)
			default:
				value = fmt.Sprint(stype)
			}
			out += fmt.Sprintf("%s: %s\n", key, value)
		}
	}
	out = strings.TrimSuffix(out, "\n")
	return []byte(out), nil
}
