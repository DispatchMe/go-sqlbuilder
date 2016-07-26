package sqlbuilder

import (
	"errors"
	"github.com/jmoiron/sqlx/reflectx"
	"reflect"
	"strings"
)

func getData(data interface{}) (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	val := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	for t.Kind() == reflect.Ptr {
		val = val.Elem()
		t = t.Elem()
	}

	if t.Kind() == reflect.Map {
		keys := val.MapKeys()

		for _, k := range keys {
			ki := k.Interface()
			if strkey, ok := ki.(string); ok {
				mp[strkey] = val.MapIndex(k).Interface()
			} else {
				return nil, errors.New("Cannot insert/update map with non-string keys")
			}
		}
		return mp, nil
	} else if t.Kind() == reflect.Struct {
		mapper := reflectx.NewMapperFunc("db", func(s string) string {
			return s
		})

		fields := mapper.FieldMap(reflect.ValueOf(data))

		// ReflectX returns nested structs. Remove any with a dot, since those are nested properties
		for fieldName, fieldValue := range fields {
			if strings.Contains(fieldName, ".") {
				continue
			}
			mp[fieldName] = fieldValue.Interface()
		}

		return mp, nil

	} else {
		return nil, errors.New("Can only insert maps and structs! (got a " + t.Kind().String() + ")")
	}

}
