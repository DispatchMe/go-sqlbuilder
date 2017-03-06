package sqlbuilder

import (
	"fmt"
	"strings"
)

type table struct {
	name            string
	joinType        int
	joinConstraints []SQLProvider
	subQuery        SQLProvider
}

func (t *table) GetSQL(cache *VarCache) string {

	if t.joinType == join_none {
		if t.subQuery != nil {
			return t.subQuery.GetSQL(cache)
		}
		return t.name
	}

	joinType := joinTypes[t.joinType]

	ons := make([]string, len(t.joinConstraints))

	for i, c := range t.joinConstraints {
		ons[i] = c.GetSQL(cache)
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
func OnColumn(field, otherField string) SQLProvider {
	return Raw{field + " = " + otherField}
}

// Just a passthru to make people feel better
func OnExpression(expr SQLProvider) SQLProvider {
	return expr
}

// Add a JOIN component. Constraints can be any number of sqlPRoviders, but you should only
// use OnColumn and OnValue
func (q *Query) join(joinType int, tableName string, constraints ...SQLProvider) *Query {
	newTable := &table{
		joinType:        joinType,
		name:            tableName,
		joinConstraints: constraints,
	}

	q.tables = append(q.tables, newTable)

	return q
}

func (q *Query) InnerJoin(tableName string, constraints ...SQLProvider) *Query {
	return q.join(join_inner, tableName, constraints...)
}

func (q *Query) LeftJoin(tableName string, constraints ...SQLProvider) *Query {
	return q.join(join_left, tableName, constraints...)
}

func (q *Query) RightJoin(tableName string, constraints ...SQLProvider) *Query {
	return q.join(join_right, tableName, constraints...)
}

func (q *Query) OuterJoin(tableName string, constraints ...SQLProvider) *Query {
	return q.join(join_outer, tableName, constraints...)
}

// Get the alias of a table that has already been joined on the query. This is useful for when you construct the query using several different components and potentially more than one component may need to join the same table. Important to note, however that this function will return the first joined table that matches. If you join the same table twice with different aliases, you should not rely on this function.
func (q *Query) GetJoinAlias(tableName string) (found bool, alias string) {
	var spl []string
	var lowerName string

	for _, t := range q.tables {
		lowerName = strings.ToLower(t.name)
		spl = strings.Split(lowerName, " as ")

		if spl[0] == tableName {
			found = true
			if len(spl) > 1 {
				alias = spl[1]
			} else {
				alias = spl[0]
			}
			break
		}
	}

	return
}

// Join a table/alias comboination if it hasn't already been joined
func (q *Query) lazyJoin(joinType int, tableName string, constraints ...SQLProvider) *Query {
	for _, t := range q.tables {
		if t.name == tableName && t.joinType == joinType {
			return q
		}
	}

	return q.join(joinType, tableName, constraints...)
}

func (q *Query) LazyInnerJoin(tableName string, constraints ...SQLProvider) *Query {
	return q.lazyJoin(join_inner, tableName, constraints...)
}

func (q *Query) LazyLeftJoin(tableName string, constraints ...SQLProvider) *Query {
	return q.lazyJoin(join_left, tableName, constraints...)
}

func (q *Query) LazyRightJoin(tableName string, constraints ...SQLProvider) *Query {
	return q.lazyJoin(join_right, tableName, constraints...)
}

func (q *Query) LazyOuterJoin(tableName string, constraints ...SQLProvider) *Query {
	return q.lazyJoin(join_outer, tableName, constraints...)
}
