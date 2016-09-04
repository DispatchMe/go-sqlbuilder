// Package sqlbuilder facilitates programamatically generating SQL queries using a chainable interface.

package sqlbuilder

import (
	"database/sql"
	"fmt"
	sqlx "github.com/jmoiron/sqlx"
	"github.com/visionmedia/go-debug"
)

var Debug = debug.Debug("sql")

const (
	gate_and = iota
	gate_or  = iota
)

const (
	join_none  = iota
	join_inner = iota
	join_left  = iota
	join_right = iota
	join_outer = iota
)

const (
	action_select = iota
	action_insert = iota
	action_update = iota
	action_delete = iota
	action_union  = iota
)

type Query struct {
	action    int
	fields    []string
	tables    []*table
	cache     *varCache
	having    *constraint
	where     *constraint
	groups    groups
	ordering  ordering
	data      map[string]interface{}
	limit     int
	offset    int
	returning string
	unions    []SQLProvider
}

type SQLProvider interface {
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

func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

func newQuery() *Query {
	q := new(Query)
	q.cache = new(varCache)
	return q
}

// Generate the SQL for this query. Returns the generated SQL (string), and a slice of arbitrary values to pass to sql.DB.Exec or sql.DB.Query
func (q *Query) GetSQL() (string, []interface{}) {
	cache := &varCache{}
	return q.getSQL(cache), cache.vars
}

// This satisfies the SQLProvider interface so we can use subqueries
func (q *Query) getSQL(cache *varCache) string {
	var sql string

	switch q.action {
	case action_select:
		sql = q.getSelectSQL(cache)
	case action_insert:
		sql = q.getInsertSQL(cache)
	case action_update:
		sql = q.getUpdateSQL(cache)
	case action_delete:
		sql = q.getDeleteSQL(cache)
	case action_union:
		sql = q.getUnionSQL(cache)
	}
	return sql
}

// Execute a write query (INSERT/UPDATE/DELETE) on a given SQL database
func (q *Query) ExecWrite(db *sqlx.DB) (sql.Result, error) {
	sql, vars := q.GetSQL()
	Debug(sql)
	return db.Exec(sql, vars...)
}

// Execute a read query (SELECT) on a given SQL database
func (q *Query) ExecRead(db *sqlx.DB) (*sqlx.Rows, error) {
	sql, vars := q.GetSQL()
	Debug(sql)
	return db.Queryx(sql, vars...)
}

func (q *Query) GetResult(db *sqlx.DB, result interface{}) error {
	sql, vars := q.GetSQL()
	Debug(sql)
	return db.Get(result, sql, vars...)
}

func (q *Query) GetValue(db *sqlx.DB, val interface{}) error {
	results, err := q.ExecRead(db)
	if err != nil {
		return err
	}

	defer results.Close()
	results.Next()
	err = results.Scan(val)
	return err
}
