package sqlbuilder

import (
	"fmt"
	"strings"
)

// Start an INSERT query
func Insert(data interface{}) *Query {
	query := newQuery()

	formatted, err := getData(data)
	if err != nil {
		// For now just panic. It's essentially a "compiler error" enforced at runtime
		// because we need to accept multiple types
		panic(err)
	}

	query.data = formatted
	query.action = ACTION_INSERT
	return query

}

// Set the table to insert into
func (q *Query) Into(tableName string) *Query {
	// Technically does the same thing as From
	return q.From(tableName)
}

func (q *Query) getInsertSQL(cache *varCache) string {
	keys := make([]string, len(q.data))

	fields := make([]string, len(q.data))

	i := 0
	for k, d := range q.data {
		keys[i] = cache.add(d)
		fields[i] = k
		i++
	}

	// Make sure there's only one table and it is JOIN_NONE
	if len(q.tables) != 1 {
		panic("Cannot run INSERT on multiple tables!")
	} else if q.tables[0].joinType != JOIN_NONE {
		panic("Cannot run INSERT on table with join type")
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", q.tables[0].name, strings.Join(fields, ", "), strings.Join(keys, ", "))
}
