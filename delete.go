package sqlbuilder

import (
	"strings"
)

func Delete() *Query {
	query := newQuery()
	query.action = ACTION_DELETE

	return query
}

func (q *Query) getDeleteSQL(cache *varCache) string {
	components := []string{"DELETE FROM"}

	common := q.getCommonQueryComponents(cache)

	order := []string{"from", "join", "where", "groupBy", "having", "orderBy", "limit", "skip"}

	for _, o := range order {
		if val, ok := common[o]; ok {
			components = append(components, val)
		}
	}

	return strings.Join(components, " ")
}
