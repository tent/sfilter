package sfilter

import (
	"fmt"
	"reflect"
	"strings"
)

func isStructSlice(v reflect.Value) bool {
	v = reflect.Indirect(v)
	if v.Kind() == reflect.Slice {
		elemType := v.Type().Elem()
		elemKind := elemType.Kind()
		if elemKind == reflect.Struct || elemKind == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct {
			return true
		}
	}
	return false
}

type marshaler interface {
	MarshalJSON() ([]byte, error)
}

// Recursively traverse struct v and return a map with values that are tagged
// with matching tags. If v is self-referential, this will
// result in an infinite loop.
func Map(v interface{}, tags ...string) (map[string]interface{}, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("sfilter: no tags provided")
	}

	src := reflect.Indirect(reflect.ValueOf(v))
	if src.Kind() != reflect.Struct {
		return nil, fmt.Errorf("sfilter: %T is not a struct or struct pointer")
	}

	srcType := src.Type()
	dest := make(map[string]interface{})
	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		fieldType := srcType.Field(i)
		fieldTag := fieldType.Tag.Get("sfilter")

		if fieldTag == "" {
			continue
		}

		fieldTags := strings.Split(fieldTag, ",")
		keep := false
	tagloop:
		for _, t := range tags {
			for _, ft := range fieldTags {
				if t == ft {
					keep = true
					break tagloop
				}
			}
		}
		if !keep {
			continue
		}

		name, options := parseTag(fieldType.Tag.Get("json"))
		if name == "" {
			name = fieldType.Name
		}

		m, ok := field.Interface().(marshaler)
		if !ok {
			// T doesn't match the interface. Check against *T too.
			if field.Kind() != reflect.Ptr && field.CanAddr() {
				m, ok = field.Addr().Interface().(marshaler)
				if ok {
					field = field.Addr()
				}
			}
		}
		if ok && (field.Kind() != reflect.Ptr || !field.IsNil()) {
			dest[name] = m
			continue
		}

		if !field.IsValid() || options.Contains("omitempty") && isEmptyValue(field) {
			continue
		}
		field = reflect.Indirect(field)
		var err error
		if _, ok := field.Interface().(marshaler); ok {
			dest[name] = field.Interface()
		} else if field.Kind() == reflect.Struct {
			dest[name], err = Map(field.Interface(), tags...)
			if err != nil {
				return nil, err
			}
		} else if isStructSlice(field) {
			slice := make([]map[string]interface{}, field.Len())
			for i := 0; i < field.Len(); i++ {
				slice[i], err = Map(field.Index(i).Interface(), tags...)
				if err != nil {
					return nil, err
				}
			}
			dest[name] = slice
		} else {
			dest[name] = field.Interface()
		}
	}

	return dest, nil
}
