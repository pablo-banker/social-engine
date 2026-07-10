package validation

import (
	"reflect"
	"strings"
)

func getJSONName(obj interface{}, field string) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	f, ok := t.FieldByName(field)
	if !ok {
		return field
	}
	tag := f.Tag.Get("json")
	if tag == "" || tag == "-" {
		return field
	}
	return strings.Split(tag, ",")[0]
}
