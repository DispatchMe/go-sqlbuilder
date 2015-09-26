package sqlbuilder

import (
	"reflect"
	"strings"
)

// "="" expression. "field = value"
type Equal struct {
	Field string
	Value interface{}
}

func getSQLFromInterface(cache *varCache, i interface{}) string {
	if provider, ok := i.(sqlProvider); ok {
		return "(" + provider.getSQL(cache) + ")"
	} else {
		return cache.add(i)
	}
}

func (e Equal) getSQL(cache *varCache) string {
	return e.Field + " = " + getSQLFromInterface(cache, e.Value)
}

// "!=" expression. "field != value"
type NotEqual struct {
	Field string
	Value interface{}
}

func (e NotEqual) getSQL(cache *varCache) string {
	return e.Field + " != " + getSQLFromInterface(cache, e.Value)
}

// ">" expression. "field > value"
type GreaterThan struct {
	Field string
	Value interface{}
}

func (e GreaterThan) getSQL(cache *varCache) string {
	return e.Field + " > " + getSQLFromInterface(cache, e.Value)
}

// "<" expression. "field < value"
type LessThan struct {
	Field string
	Value interface{}
}

func (e LessThan) getSQL(cache *varCache) string {
	return e.Field + " < " + getSQLFromInterface(cache, e.Value)
}

// ">=" expression. "field >= value"
type GreaterOrEqual struct {
	Field string
	Value interface{}
}

func (e GreaterOrEqual) getSQL(cache *varCache) string {
	return e.Field + " >= " + getSQLFromInterface(cache, e.Value)
}

// "<=" expression. "field <= value"
type LessOrEqual struct {
	Field string
	Value interface{}
}

func (e LessOrEqual) getSQL(cache *varCache) string {
	return e.Field + " <= " + getSQLFromInterface(cache, e.Value)
}

func getInKeys(val interface{}, cache *varCache) string {
	if provider, ok := val.(sqlProvider); ok {
		return provider.getSQL(cache)
	}
	// Make sure value is a slice
	t := reflect.TypeOf(val)
	v := reflect.ValueOf(val)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Slice {
		panic(`Value for "In" expression must be a slice`)
	}

	keys := make([]string, v.Len())

	for i := 0; i < v.Len(); i++ {
		keys[i] = cache.add(v.Index(i).Interface())
	}
	return strings.Join(keys, ", ")

}

// "IN" expression. "field IN (values...)". Values must be a slice, otherwise this function will panic.
type In struct {
	Field string
	Value interface{}
}

func (e In) getSQL(cache *varCache) string {
	return e.Field + " IN (" + getInKeys(e.Value, cache) + ")"
}

// "NOT IN" expression. "field NOT IN (values...)". Values must be a slice, otherwise this function will panic.
type NotIn struct {
	Field string
	Value interface{}
}

func (e NotIn) getSQL(cache *varCache) string {
	return e.Field + " NOT IN (" + getInKeys(e.Value, cache) + ")"
}

// "LIKE" expression. "field LIKE value"
type Like struct {
	Field string
	Value interface{}
}

func (e Like) getSQL(cache *varCache) string {
	return e.Field + " LIKE " + cache.add(e.Value)
}

// "IS NULL" expression. "field IS NULL"
type IsNull struct {
	Field string
}

func (e IsNull) getSQL(cache *varCache) string {
	return e.Field + " IS NULL"
}

// "IS NOT NULL" expression. "field IS NOT NULL"
type IsNotNull struct {
	Field string
}

func (e IsNotNull) getSQL(cache *varCache) string {
	return e.Field + " IS NOT NULL"
}

// Raw expression. BE CAREFUL!!! THERE IS NO SQL INJECTION PROTECTION HERE. This is used internally for joining on columns. Use at your own risk
type Raw struct {
	Expr string
}

func (e Raw) getSQL(cache *varCache) string {
	return e.Expr
}
