package sqlbuilder

import (
	"fmt"
)

const (
	GATE_AND = iota
	GATE_OR  = iota
)

const (
	JOIN_NONE  = iota
	JOIN_INNER = iota
	JOIN_LEFT  = iota
	JOIN_RIGHT = iota
	JOIN_OUTER = iota
)

const (
	ACTION_SELECT = iota
	ACTION_INSERT = iota
	ACTION_UPDATE = iota
	ACTION_DELETE = iota
)

type Query struct {
	action   int
	fields   []string
	tables   []*table
	cache    *varCache
	having   *constraint
	where    *constraint
	groups   groups
	ordering ordering
	data     map[string]interface{}
	limit    int
	skip     int
}

type sqlProvider interface {
	getSQL(cache *varCache) string
}

type varCache struct {
	vars []interface{}
}

func (v *varCache) add(val interface{}) string {
	v.vars = append(v.vars, val)
	return fmt.Sprintf("$%d", len(v.vars))
}

type group struct {
	field      string
	descending bool
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

func (q *Query) Skip(skip int) *Query {
	q.skip = skip
	return q
}

func newQuery() *Query {
	q := new(Query)
	q.cache = new(varCache)
	return q
}

func (q *Query) GetSQL() (string, []interface{}) {
	cache := &varCache{}

	var sql string
	switch q.action {
	case ACTION_SELECT:
		sql = q.getSelectSQL(cache)
	case ACTION_INSERT:
		sql = q.getInsertSQL(cache)
	case ACTION_UPDATE:
		sql = q.getUpdateSQL(cache)
	case ACTION_DELETE:
		sql = q.getDeleteSQL(cache)
	}

	return sql, cache.vars

}
