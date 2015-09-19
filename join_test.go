package sqlbuilder

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJoins(t *testing.T) {
	Convey("Join collection SQL generation", t, func() {
		cache := &varCache{}
		Convey("single constraint", func() {
			query := Select("*")

			query.Join(JOIN_INNER, "foos AS foobees", OnColumn("bars.foo_id", "=", "foos.id"))

			So(len(query.tables), ShouldEqual, 1)
			So(query.tables[0].getSQL(cache), ShouldEqual, `INNER JOIN foos AS foobees ON bars.foo_id = foos.id`)
		})

		Convey("multiple constraints", func() {
			query := Select("*")

			query.Join(JOIN_INNER, "foos AS foobees", OnColumn("bars.foo_id", "=", "foos.id"), OnValue("foos.doop", "=", "foop"))

			So(len(query.tables), ShouldEqual, 1)
			So(query.tables[0].getSQL(cache), ShouldEqual, `INNER JOIN foos AS foobees ON bars.foo_id = foos.id AND ON foos.doop = $1`)
			So(cache.vars[0], ShouldEqual, "foop")
		})
	})
}
