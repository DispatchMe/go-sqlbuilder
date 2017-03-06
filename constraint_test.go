package sqlbuilder

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConstraints(t *testing.T) {
	Convey("Constraints SQL generation", t, func() {
		cache := &VarCache{}

		Convey("and combined", func() {
			c := new(constraint)
			c.gate = gate_and

			c.addChild(Equal{"foo", 10})
			c.addChild(Equal{"bar", "bar"})
			sql := c.GetSQL(cache)
			So(sql, ShouldEqual, `(foo = $1 AND bar = $2)`)
			So(len(cache.vars), ShouldEqual, 2)
			So(cache.vars[0], ShouldEqual, 10)
			So(cache.vars[1], ShouldEqual, "bar")

		})

		Convey("or combined", func() {
			c := new(constraint)
			c.gate = gate_or
			c.addChild(Equal{"foo", 10})
			c.addChild(Equal{"bar", "bar"})

			sql := c.GetSQL(cache)
			So(sql, ShouldEqual, `(foo = $1 OR bar = $2)`)
			So(len(cache.vars), ShouldEqual, 2)
			So(cache.vars[0], ShouldEqual, 10)
			So(cache.vars[1], ShouldEqual, "bar")
		})

		Convey("complex", func() {
			c := new(constraint)
			c.gate = gate_or

			c.children = []SQLProvider{
				Equal{"foo", 10},
				&constraint{
					gate: gate_and,
					children: []SQLProvider{
						Equal{"bar", "bar"},
						Equal{"baz", "baz"},
					},
				},
			}

			sql := c.GetSQL(cache)
			So(sql, ShouldEqual, `(foo = $1 OR (bar = $2 AND baz = $3))`)
			So(len(cache.vars), ShouldEqual, 3)
			So(cache.vars[0], ShouldEqual, 10)
			So(cache.vars[1], ShouldEqual, "bar")
			So(cache.vars[2], ShouldEqual, "baz")
		})
	})
}
