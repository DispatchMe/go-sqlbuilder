package sqlbuilder

import (
	"strings"
)

// Run an UPDATE query
func Update(tableName string) *Query {
	query := newQuery()
	query.action = action_update
	// Technically does the same thing as From
	return query.From(tableName)
}

// Set the data for an update. This can either be a struct or a map[string]interface{}.
// Structs can use the `db` tag to designate alternate column names
func (q *Query) Set(data interface{}) *Query {
	formatted, err := getData(data)
	if err != nil {
		// For now just panic. It's essentially a "compiler error" enforced at runtime
		// because we need to accept multiple types
		panic(err)
	}

	q.data = formatted

	return q
}

func (q *Query) getUpdateSQL(cache *varCache) string {

	clauses := make([]string, len(q.data))
	i := 0

	for key, val := range q.data {
		clauses[i] = key + "=" + cache.add(val)
		i++
	}

	components := []string{"UPDATE"}

	common := q.getCommonQueryComponents(cache)

	common["set"] = "SET " + strings.Join(clauses, ", ")

	order := []string{"from", "join", "set", "where", "groupBy", "having", "orderBy", "limit", "offset"}

	for _, o := range order {
		if val, ok := common[o]; ok {
			components = append(components, val)
		}
	}

	return strings.Join(components, " ")
}
