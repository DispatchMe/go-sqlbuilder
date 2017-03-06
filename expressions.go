package sqlbuilder

import (
	"fmt"
	"reflect"
	"strings"
)

// "="" expression. "field = value"
type Equal struct {
	Field string
	Value interface{}
}

func GetSQLFromInterface(cache *VarCache, i interface{}) string {
	if provider, ok := i.(SQLProvider); ok {
		return "(" + provider.GetSQL(cache) + ")"
	} else {
		return cache.add(i)
	}
}

func (e Equal) GetSQL(cache *VarCache) string {
	return e.Field + " = " + GetSQLFromInterface(cache, e.Value)
}

// "!=" expression. "field != value"
type NotEqual struct {
	Field string
	Value interface{}
}

func (e NotEqual) GetSQL(cache *VarCache) string {
	return e.Field + " != " + GetSQLFromInterface(cache, e.Value)
}

// ">" expression. "field > value"
type GreaterThan struct {
	Field string
	Value interface{}
}

func (e GreaterThan) GetSQL(cache *VarCache) string {
	return e.Field + " > " + GetSQLFromInterface(cache, e.Value)
}

// "<" expression. "field < value"
type LessThan struct {
	Field string
	Value interface{}
}

func (e LessThan) GetSQL(cache *VarCache) string {
	return e.Field + " < " + GetSQLFromInterface(cache, e.Value)
}

// ">=" expression. "field >= value"
type GreaterOrEqual struct {
	Field string
	Value interface{}
}

func (e GreaterOrEqual) GetSQL(cache *VarCache) string {
	return e.Field + " >= " + GetSQLFromInterface(cache, e.Value)
}

// "<=" expression. "field <= value"
type LessOrEqual struct {
	Field string
	Value interface{}
}

func (e LessOrEqual) GetSQL(cache *VarCache) string {
	return e.Field + " <= " + GetSQLFromInterface(cache, e.Value)
}

func getInKeys(val interface{}, cache *VarCache) string {
	if provider, ok := val.(SQLProvider); ok {
		return provider.GetSQL(cache)
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

// "IN" expression. "field IN (values...)". Values must be a slice or a subquery, otherwise this function will panic.
type In struct {
	Field string
	Value interface{}
}

func (e In) GetSQL(cache *VarCache) string {
	return e.Field + " IN (" + getInKeys(e.Value, cache) + ")"
}

// "NOT IN" expression. "field NOT IN (values...)". Values must be a slice or a subquery, otherwise this function will panic.
type NotIn struct {
	Field string
	Value interface{}
}

func (e NotIn) GetSQL(cache *VarCache) string {
	return e.Field + " NOT IN (" + getInKeys(e.Value, cache) + ")"
}

// "LIKE" expression. "field LIKE value"
type Like struct {
	Field string
	Value interface{}
}

func (e Like) GetSQL(cache *VarCache) string {
	return e.Field + " LIKE " + cache.add(e.Value)
}

// "IS NULL" expression. "field IS NULL"
type IsNull struct {
	Field string
}

func (e IsNull) GetSQL(cache *VarCache) string {
	return e.Field + " IS NULL"
}

// "IS NOT NULL" expression. "field IS NOT NULL"
type IsNotNull struct {
	Field string
}

func (e IsNotNull) GetSQL(cache *VarCache) string {
	return e.Field + " IS NOT NULL"
}

// Raw expression. BE CAREFUL!!! THERE IS NO SQL INJECTION PROTECTION HERE. This is used internally for joining on columns. Use at your own risk
type Raw struct {
	Expr string
}

func (e Raw) GetSQL(cache *VarCache) string {
	return e.Expr
}

type PGOverlap struct {
	Field string
	Value []string
}

func (e PGOverlap) GetSQL(cache *VarCache) string {
	parsed := make([]string, len(e.Value))

	for idx, item := range e.Value {
		parsed[idx] = cache.add(item)
	}
	return fmt.Sprintf("%s && ARRAY[%s]::varchar[]", e.Field, strings.Join(parsed, ","))
}

type Expression struct {
	Format string
	Values []interface{}
}

func (e Expression) GetSQL(cache *VarCache) string {
	parsed := make([]interface{}, len(e.Values))

	for idx, item := range e.Values {
		parsed[idx] = cache.add(item)
	}

	return fmt.Sprintf(e.Format, parsed...)
}
