package sqlbuilder

import (
	"errors"
	"github.com/jmoiron/sqlx/reflectx"
	"reflect"
)

func getData(data interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	for t.Kind() == reflect.Ptr {
		val = val.Elem()
		t = t.Elem()
	}

	if t.Kind() == reflect.Map {
		d, ok := data.(map[string]interface{})
		if !ok {
			return nil, errors.New("Map must be a map[string]interface{}")
		}
		return d, nil
	} else if t.Kind() == reflect.Struct {
		mp := make(map[string]interface{})
		mapper := reflectx.NewMapperFunc("db", func(s string) string {
			return s
		})

		fields := mapper.TypeMap(reflect.TypeOf(data))
		idx := fields.Index
		for _, x := range idx {
			mp[x.Name] = mapper.FieldByName(val, x.Name).Interface()
		}

		return mp, nil

	} else {
		return nil, errors.New("Can only insert maps and structs! (got a " + t.Kind().String() + ")")
	}

}
