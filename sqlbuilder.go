// Package sqlbuilder facilitates programamatically generating SQL queries using a chainable interface.

package sqlbuilder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	sqlx "github.com/jmoiron/sqlx"
	"github.com/visionmedia/go-debug"
	"os"
	"strings"
)

var debugEnabled = strings.Contains(os.Getenv("DEBUG"), "sql")

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
	action_count  = iota
)

type PlaceholderFunction func(index int) string

type Query struct {
	action      int
	fields      []string
	tables      []*table
	cache       *VarCache
	having      *constraint
	where       *constraint
	groups      groups
	ordering    ordering
	data        map[string]interface{}
	limit       int
	offset      int
	placeholder PlaceholderFunction
	returning   string
	unions      []SQLProvider
}

type SQLProvider interface {
	GetSQL(cache *VarCache) string
}

type VarCache struct {
	placeholder PlaceholderFunction
	vars        []interface{}
}

func (v *VarCache) add(val interface{}) string {
	v.vars = append(v.vars, val)
	if v.placeholder != nil {
		return v.placeholder(len(v.vars))
	}
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

// Change the prepared statement placeholder (the question mark in this example) (INSERT INTO _ (?, ?, ?) VALUES())
func (q *Query) Placeholder(placeholder PlaceholderFunction) *Query {
	q.placeholder = placeholder
	return q
}

func newQuery() *Query {
	q := new(Query)
	q.cache = new(VarCache)
	return q
}

// Generate the SQL for this query. Returns the generated SQL (string), and a slice of arbitrary values to pass to sql.DB.Exec or sql.DB.Query
func (q *Query) GetFullSQL() (string, []interface{}) {
	cache := &VarCache{
		placeholder: q.placeholder,
	}
	return q.GetSQL(cache), cache.vars
}

// This satisfies the SQLProvider interface so we can use subqueries
func (q *Query) GetSQL(cache *VarCache) string {
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
	case action_count:
		sql = q.getCountSQL(cache)
	}
	return sql
}

func (q *Query) GetCount(db *sqlx.DB) (int, error) {
	var count int
	prevAction := q.action
	q.action = action_count
	defer func() {
		q.action = prevAction
	}()
	err := q.GetValue(db, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Execute a write query (INSERT/UPDATE/DELETE) on a given SQL database
func (q *Query) ExecWrite(db *sqlx.DB) (sql.Result, error) {
	sql, vars := q.GetFullSQL()

	if debugEnabled {
		marshaled, _ := json.Marshal(vars)
		Debug("%s, %s", sql, string(marshaled))
	}

	return db.Exec(sql, vars...)
}

// Execute a read query (SELECT) on a given SQL database
func (q *Query) ExecRead(db *sqlx.DB) (*sqlx.Rows, error) {
	sql, vars := q.GetFullSQL()

	if debugEnabled {
		marshaled, _ := json.Marshal(vars)
		Debug("%s, %s", sql, string(marshaled))
	}
	return db.Queryx(sql, vars...)
}

func (q *Query) GetResult(db *sqlx.DB, result interface{}) error {
	sql, vars := q.GetFullSQL()

	if debugEnabled {
		marshaled, _ := json.Marshal(vars)
		Debug("%s, %s", sql, string(marshaled))
	}

	return db.Get(result, sql, vars...)
}

func (q *Query) GetValue(db *sqlx.DB, val interface{}) error {
	results, err := q.ExecRead(db)
	if err != nil {
		return err
	}

	defer results.Close()
	results.Next()
	err = results.Err()
	if err != nil {
		return err
	}

	err = results.Scan(val)
	return err
}
