package sqlbuilder

import (
	"strings"
)

// Run a DELETE query
func Delete() *Query {
	query := newQuery()
	query.action = action_delete

	return query
}

func (q *Query) getDeleteSQL(cache *varCache) string {
	components := []string{"DELETE FROM"}

	common := q.getCommonQueryComponents(cache)

	order := []string{"from", "join", "where", "groupBy", "having", "orderBy", "limit", "offset"}

	for _, o := range order {
		if val, ok := common[o]; ok {
			components = append(components, val)
		}
	}

	return strings.Join(components, " ")
}
