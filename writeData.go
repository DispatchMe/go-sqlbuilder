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
		addStructFieldsToMap(structs.Fields(data), mp)
		return mp, nil

	} else {
		return nil, errors.New("Can only insert maps and structs! (got a " + t.Kind().String() + ")")
	}
}

func addStructFieldsToMap(fields []*structs.Field, mp map[string]interface{}) {
	for _, f := range fields {
		if f.IsEmbedded() {
			addStructFieldsToMap(f.Fields(), mp)
		} else {
			addStructFieldToMap(f, mp)
		}
	}
}

func addStructFieldToMap(field *structs.Field, mp map[string]interface{}) {
	tag := field.Tag("db")
	if strings.Contains(tag, ",") {
		spl := strings.Split(tag, ",")
		tag = spl[0]
		if len(spl) > 1 && spl[1] == "omitempty" && field.IsZero() {
			return
		}
	}

	if tag == "" {
		tag = field.Name()
	}

	mp[tag] = field.Value()
}
