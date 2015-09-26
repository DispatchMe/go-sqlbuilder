package sqlbuilder

import (
	"errors"
	"github.com/jmoiron/sqlx/reflectx"
	"reflect"
)

// This is mostly used for testing, but it's nice to see the SQL output in order of input. This wraps a map[string]interface and maintains the order of key/value insertion when looping through with hasNext/getNext.
type orderedMap struct {
	keys []string
	data map[string]interface{}

	idx int
}

func (o *orderedMap) set(key string, val interface{}) {
	o.keys = append(o.keys, key)
	if o.data == nil {
		o.data = make(map[string]interface{})
	}
	o.data[key] = val
}

func (o *orderedMap) hasNext() bool {
	return o.idx < len(o.keys)
}

func (o *orderedMap) getNext() (key string, value interface{}) {
	key = o.keys[o.idx]
	value = o.data[key]
	o.idx++
	return
}

func (o *orderedMap) rewind() {
	o.idx = 0
}

func getData(data interface{}) (*orderedMap, error) {
	omap := &orderedMap{}

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
				omap.set(strkey, val.MapIndex(k).Interface())
			} else {
				return nil, errors.New("Cannot insert/update map with non-string keys")
			}
		}
		return omap, nil
	} else if t.Kind() == reflect.Struct {
		mapper := reflectx.NewMapperFunc("db", func(s string) string {
			return s
		})

		fields := mapper.TypeMap(reflect.TypeOf(data))
		idx := fields.Index
		for _, x := range idx {
			omap.set(x.Name, mapper.FieldByName(val, x.Name).Interface())
		}

		return omap, nil

	} else {
		return nil, errors.New("Can only insert maps and structs! (got a " + t.Kind().String() + ")")
	}

}
