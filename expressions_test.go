package sqlbuilder

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type exprParam struct {
	name     string
	expr     SQLProvider
	expected string
	vars     []interface{}
}

func TestExpressions(t *testing.T) {
	params := []exprParam{
		{"Equal", Equal{"foo", "bar"}, "foo = $1", []interface{}{"bar"}},
		{"NotEqual", NotEqual{"foo", 10}, "foo != $1", []interface{}{10}},
		{"GreaterThan", GreaterThan{"age", 5}, "age > $1", []interface{}{5}},
		{"LessThan", LessThan{"age", 5}, "age < $1", []interface{}{5}},
		{"GreaterOrEqual", GreaterOrEqual{"age", 5}, "age >= $1", []interface{}{5}},
		{"LessOrEqual", LessOrEqual{"age", 5}, "age <= $1", []interface{}{5}},
		{"In", In{"category", []string{"books", "movies", "music"}}, "category IN ($1, $2, $3)", []interface{}{"books", "movies", "music"}},
		{"NotIn", NotIn{"category", []string{"books", "movies", "music"}}, "category NOT IN ($1, $2, $3)", []interface{}{"books", "movies", "music"}},
		{"Like", Like{"title", "%america%"}, "title LIKE $1", []interface{}{"%america%"}},
		{"IsNull", IsNull{"author_id"}, "author_id IS NULL", nil},
		{"IsNotNull", IsNotNull{"author_id"}, "author_id IS NOT NULL", nil},
		{"In subquery", In{"game_id", Select("id").From("games").Where(Equal{"type", "football"})}, "game_id IN (SELECT id FROM games WHERE type = $1)", []interface{}{"football"}},
		{"Equal to subquery", Equal{"sum", Select("SUM(count)").From("foos")}, "sum = (SELECT SUM(count) FROM foos)", nil},
	}
	Convey("Expressions", t, func() {
		cache := &VarCache{}
		for _, p := range params {
			Convey(p.name, func() {
				sql := p.expr.GetSQL(cache)
				So(sql, ShouldEqual, p.expected)

				for i, v := range cache.vars {
					So(p.vars[i], ShouldEqual, v)
				}
			})
		}
	})
}
