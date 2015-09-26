package sqlbuilder

import (
	"fmt"
	"strings"
)

type table struct {
	name            string
	joinType        int
	joinConstraints []sqlProvider
	subQuery        sqlProvider
}

func (t *table) getSQL(cache *varCache) string {

	if t.joinType == join_none {
		if t.subQuery != nil {
			return t.subQuery.getSQL(cache)
		}
		return t.name
	}

	joinType := joinTypes[t.joinType]

	ons := make([]string, len(t.joinConstraints))

	for i, c := range t.joinConstraints {
		ons[i] = c.getSQL(cache)
	}

	return fmt.Sprintf("%s JOIN %s ON %s", joinType, t.name, strings.Join(ons, " AND "))
}

var joinTypes = map[int]string{
	join_inner: "INNER",
	join_left:  "LEFT",
	join_right: "RIGHT",
	join_outer: "OUTER",
}

// Used within a JOIN context, joins on a particular column link
func OnColumn(field, otherField string) sqlProvider {
	return Raw{field + " = " + otherField}
}

// Just a passthru to make people feel better
func OnExpression(expr sqlProvider) sqlProvider {
	return expr
}

// Add a JOIN component. Constraints can be any number of sqlPRoviders, but you should only
// use OnColumn and OnValue
func (q *Query) join(joinType int, tableName string, constraints ...sqlProvider) *Query {
	newTable := &table{
		joinType:        joinType,
		name:            tableName,
		joinConstraints: constraints,
	}

	q.tables = append(q.tables, newTable)

	return q
}

func (q *Query) InnerJoin(tableName string, constraints ...sqlProvider) *Query {
	return q.join(join_inner, tableName, constraints...)
}

func (q *Query) LeftJoin(tableName string, constraints ...sqlProvider) *Query {
	return q.join(join_left, tableName, constraints...)
}

func (q *Query) RightJoin(tableName string, constraints ...sqlProvider) *Query {
	return q.join(join_right, tableName, constraints...)
}

func (q *Query) OuterJoin(tableName string, constraints ...sqlProvider) *Query {
	return q.join(join_outer, tableName, constraints...)
}
