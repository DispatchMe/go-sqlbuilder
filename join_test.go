package sqlbuilder

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJoins(t *testing.T) {
	Convey("Join collection SQL generation", t, func() {
		cache := &VarCache{}
		Convey("single constraint", func() {
			query := Select("*")

			query.InnerJoin("foos AS foobees", OnColumn("bars.foo_id", "foos.id"))

			So(len(query.tables), ShouldEqual, 1)
			So(query.tables[0].GetSQL(cache), ShouldEqual, `INNER JOIN foos AS foobees ON bars.foo_id = foos.id`)
		})

		Convey("multiple constraints", func() {
			query := Select("*")

			query.LeftJoin("foos AS foobees", OnColumn("bars.foo_id", "foos.id"), OnExpression(Equal{"foos.doop", "foop"}))

			So(len(query.tables), ShouldEqual, 1)
			So(query.tables[0].GetSQL(cache), ShouldEqual, `LEFT JOIN foos AS foobees ON bars.foo_id = foos.id AND foos.doop = $1`)
			So(cache.vars[0], ShouldEqual, "foop")
		})

		Convey("Lazy join", func() {
			query := Select("*")

			query.LazyInnerJoin("foos AS foobees", OnColumn("bars.foo_id", "foos.id"))
			query.LazyInnerJoin("foos AS foobees", OnColumn("bars.foo_id", "foos.id"))

			So(len(query.tables), ShouldEqual, 1)
		})
	})
}
