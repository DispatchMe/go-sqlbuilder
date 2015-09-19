package sqlbuilder

import (
	"fmt"
	"strings"
)

type table struct {
	name            string
	joinType        int
	joinConstraints []sqlProvider
}

func (t *table) getSQL(cache *varCache) string {

	if t.joinType == JOIN_NONE {
		return t.name
	}

	joinType := joinTypes[t.joinType]

	ons := make([]string, len(t.joinConstraints))

	for i, c := range t.joinConstraints {
		ons[i] = c.getSQL(cache)
	}

	return fmt.Sprintf("%s JOIN %s ON %s", joinType, t.name, strings.Join(ons, " AND ON "))
}

var joinTypes = map[int]string{
	JOIN_INNER: "INNER",
	JOIN_LEFT:  "LEFT",
	JOIN_RIGHT: "RIGHT",
	JOIN_OUTER: "OUTER",
}

// Used within a JOIN context, joins on a particular column link
func OnColumn(field, comparator, otherField string) *constraint {
	return &constraint{
		field:           field,
		comparator:      comparator,
		value:           otherField,
		bypassStatement: true,
	}
}

// Used within a JOIN context, joins on a particular value, IE in AND ON
func OnValue(field, comparator string, value interface{}) *constraint {
	return &constraint{
		field:      field,
		comparator: comparator,
		value:      value,
	}
}

// Add a JOIN component. Constraints can be any number of sqlPRoviders, but you should only
// use OnColumn and OnValue
func (q *Query) Join(joinType int, tableName string, constraints ...sqlProvider) *Query {
	newTable := &table{
		joinType:        joinType,
		name:            tableName,
		joinConstraints: constraints,
	}

	q.tables = append(q.tables, newTable)

	return q
}
