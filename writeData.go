package sqlbuilder

import (
	"errors"
	"github.com/fatih/structs"
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
		fields := structs.Fields(data)

		for _, f := range fields {
			tag := f.Tag("db")
			if strings.Contains(tag, ",") {
				tag = strings.Split(tag, ",")[0]
			}

			if tag == "" {
				tag = f.Name()
			}

			mp[tag] = f.Value()

		}
		return mp, nil

	} else {
		return nil, errors.New("Can only insert maps and structs! (got a " + t.Kind().String() + ")")
	}

}
