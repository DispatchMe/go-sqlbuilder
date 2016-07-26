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
	query.action = action_insert
	return query

}

// Set the table to insert into. Can be just the table name, or the table name plus the alias, e.g. `Into("foos AS bars")`
func (q *Query) Into(tableName string) *Query {
	// Technically does the same thing as From
	return q.From(tableName)
}

func (q *Query) Returning(expr string) *Query {
	q.returning = expr
	return q
}

func (q *Query) getInsertSQL(cache *varCache) string {
	keys := make([]string, len(q.data))

	fields := make([]string, len(q.data))

	i := 0

	for key, val := range q.data {
		keys[i] = cache.add(val)
		fields[i] = key
		i++
	}

	// Make sure there's only one table and it is join_none
	if len(q.tables) != 1 {
		panic("Cannot run INSERT on multiple tables!")
	} else if q.tables[0].joinType != join_none {
		panic("Cannot run INSERT on table with join type")
	}

	suffix := ""

	if q.returning != "" {
		suffix = " RETURNING " + q.returning
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", q.tables[0].name, strings.Join(fields, ", "), strings.Join(keys, ", ")) + suffix
}
