package sqlbuilder

import (
	"fmt"
	"strings"
)

const ASC = true
const DESC = false

type groups []string

func (g groups) GetSQL(cache *VarCache) string {
	return "GROUP BY " + strings.Join(g, ", ")
}

type ordering []*order

func (o ordering) GetSQL(cache *VarCache) string {
	sql := make([]string, len(o))

	var dir string
	for i, s := range o {
		if s.asc {
			dir = "ASC"
		} else {
			dir = "DESC"
		}
		sql[i] = s.expression + " " + dir
	}

	return "ORDER BY " + strings.Join(sql, ", ")
}

type order struct {
	expression string
	asc        bool
}

func (q *Query) getCountSQL(cache *VarCache) string {
	components := []string{"SELECT COUNT(*) FROM"}
	common := q.getCommonQueryComponents(cache)

	order := []string{"from", "join", "where", "groupBy", "having"}

	for _, o := range order {
		if val, ok := common[o]; ok {
			components = append(components, val)
		}
	}

	return strings.Join(components, " ")
}

func (q *Query) getSelectSQL(cache *VarCache) string {
	components := []string{"SELECT " + strings.Join(q.fields, ", ") + " FROM"}
	common := q.getCommonQueryComponents(cache)

	order := []string{"from", "join", "where", "groupBy", "having", "orderBy", "limit", "offset"}

	for _, o := range order {
		if val, ok := common[o]; ok {
			components = append(components, val)
		}
	}

	return strings.Join(components, " ")
}

func (q *Query) getCommonQueryComponents(cache *VarCache) map[string]string {
	mp := make(map[string]string)

	fromTables := []string{}
	joinTables := []string{}

	for _, t := range q.tables {
		if t.joinType == join_none {
			fromTables = append(fromTables, t.GetSQL(cache))
		} else {
			joinTables = append(joinTables, t.GetSQL(cache))
		}
	}

	mp["from"] = strings.Join(fromTables, ", ")

	if len(joinTables) > 0 {
		mp["join"] = strings.Join(joinTables, " ")
	}

	// Where?
	if q.where != nil {
		mp["where"] = "WHERE " + q.where.GetSQL(cache)
	}

	// Group?
	if len(q.groups) > 0 {
		mp["groupBy"] = q.groups.GetSQL(cache)
	}

	// Having?
	if q.having != nil {
		mp["having"] = "HAVING " + q.having.GetSQL(cache)
	}

	// Order?
	if len(q.ordering) > 0 {
		mp["orderBy"] = q.ordering.GetSQL(cache)
	}

	// Limit?
	if q.limit > 0 {
		mp["limit"] = fmt.Sprintf("LIMIT %d", q.limit)
	}

	// Offset?
	if q.offset > 0 {
		mp["offset"] = fmt.Sprintf("OFFSET %d", q.offset)
	}

	// Returning?
	if q.returning != "" {
		mp["returning"] = fmt.Sprintf("RETURNING %s", q.returning)
	}

	return mp
}

// Run a SELECT query
func Select(fields ...string) *Query {
	query := newQuery()
	query.action = action_select

	return query.Select(fields...)
}

// Select additional fields
func (q *Query) Select(fields ...string) *Query {
	q.fields = append(q.fields, fields...)
	return q
}

// In case you got a query from another package or component but want
// to change the selected fields
func (q *Query) ResetFields() *Query {
	q.fields = []string{}
	return q
}

// Alias a subquery with a certain name. Useful when you want to do something like SELECT a.column FROM (SELECT ....) a
func Alias(subquery *Query, name string) SQLProvider {
	return &alias{subquery, name}
}

type alias struct {
	query SQLProvider
	name  string
}

func (a *alias) GetSQL(cache *VarCache) string {
	return "(" + a.query.GetSQL(cache) + ") " + a.name
}

// Add a table to SELECT from. Run this multiple times for multiple tables
func (q *Query) From(tableOrQuery interface{}) *Query {
	if tableName, ok := tableOrQuery.(string); ok {
		q.tables = append(q.tables, &table{
			name: tableName,
		})
	} else if provider, ok := tableOrQuery.(SQLProvider); ok {
		q.tables = append(q.tables, &table{
			subQuery: provider,
		})
	} else {
		panic("From must be a table name (string) or a subquery!")
	}

	return q
}

// Generate a UNION clause of two or more subqueries. This can in turn be used anywhere a subquery can be used, like in FROM or IN clauses
func Union(subQuery ...SQLProvider) *Query {
	q := newQuery()
	q.unions = append(q.unions, subQuery...)
	q.action = action_union
	return q
}

func (q *Query) getUnionSQL(cache *VarCache) string {
	queries := make([]string, len(q.unions))

	for i, u := range q.unions {
		queries[i] = "(" + u.GetSQL(cache) + ")"
	}

	return strings.Join(queries, " UNION ")
}

// Add a GROUP BY clause. Use this multiple times for multiple levels of grouping
func (q *Query) GroupBy(expression string) *Query {
	q.groups = append(q.groups, expression)
	return q
}

// Add an ORDER BY clause. Use this multiple times for multiple levels of ordering
func (q *Query) OrderBy(expression string, asc bool) *Query {
	q.ordering = append(q.ordering, &order{expression, asc})
	return q
}
