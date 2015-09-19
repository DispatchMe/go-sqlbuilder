package sqlbuilder

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConstraints(t *testing.T) {
	Convey("Constraints SQL generation", t, func() {
		cache := &varCache{}
		Convey("basic usage", func() {
			c := &constraint{
				field:      "foo",
				comparator: "=",
				value:      10,
			}

			sql := c.getSQL(cache)
			So(sql, ShouldEqual, "foo = $1")
			So(len(cache.vars), ShouldEqual, 1)
			So(cache.vars[0], ShouldEqual, 10)
		})

		Convey("and combined", func() {
			c := new(constraint)
			c.gate = GATE_AND

			c.addChild(&constraint{
				field:      "foo",
				comparator: "=",
				value:      10,
			})
			c.addChild(&constraint{
				field:      "bar",
				comparator: "=",
				value:      "bar",
			})

			sql := c.getSQL(cache)
			So(sql, ShouldEqual, `(foo = $1 AND bar = $2)`)
			So(len(cache.vars), ShouldEqual, 2)
			So(cache.vars[0], ShouldEqual, 10)
			So(cache.vars[1], ShouldEqual, "bar")

		})

		Convey("or combined", func() {
			c := new(constraint)
			c.gate = GATE_OR
			c.addChild(&constraint{
				field:      "foo",
				comparator: "=",
				value:      10,
			})
			c.addChild(&constraint{
				field:      "bar",
				comparator: "=",
				value:      "bar",
			})

			sql := c.getSQL(cache)
			So(sql, ShouldEqual, `(foo = $1 OR bar = $2)`)
			So(len(cache.vars), ShouldEqual, 2)
			So(cache.vars[0], ShouldEqual, 10)
			So(cache.vars[1], ShouldEqual, "bar")
		})

		Convey("complex", func() {
			c := new(constraint)
			c.gate = GATE_OR

			c.children = []sqlProvider{
				&constraint{
					field:      "foo",
					comparator: "=",
					value:      10,
				},
				&constraint{
					gate: GATE_AND,
					children: []sqlProvider{
						&constraint{
							field:      "bar",
							comparator: "=",
							value:      "bar",
						},
						&constraint{
							field:      "baz",
							comparator: "=",
							value:      "baz",
						},
					},
				},
			}

			sql := c.getSQL(cache)
			So(sql, ShouldEqual, `(foo = $1 OR (bar = $2 AND baz = $3))`)
			So(len(cache.vars), ShouldEqual, 3)
			So(cache.vars[0], ShouldEqual, 10)
			So(cache.vars[1], ShouldEqual, "bar")
			So(cache.vars[2], ShouldEqual, "baz")
		})
	})
}
