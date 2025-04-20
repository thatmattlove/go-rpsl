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

type WithRPSL interface {
	RPSL() string
}

func RPSL(s any) (string, error) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", ErrMustBeStruct
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
			return "", fmt.Errorf("field '%s' has an invalid rpsl tag", field.Name)
		}
		value := ""
		switch stype := sval.(type) {
		case WithRPSL:
			out += stype.RPSL() + "\n"
			continue
		case uint32:
			value = strconv.FormatUint(uint64(stype), 10)
		default:
			value = fmt.Sprint(stype)
		}
		out += fmt.Sprintf("%s: %s\n", key, value)
	}
	out = strings.TrimSuffix(out, "\n")
	return out, nil
}
