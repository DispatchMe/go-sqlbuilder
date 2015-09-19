package sqlbuilder

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInsert(t *testing.T) {
	data := map[string]interface{}{
		"foo": "bar",
		"baz": 10,
	}

	Convey("insert", t, func() {
		query, vars := Insert(data).Into("people").GetSQL()

		So(len(vars), ShouldEqual, 2)
		// Problem here is map output is random, so we need to allow for both possibilities
		if query == "INSERT INTO people (foo, baz) VALUES ($1, $2)" {
			So(vars[0], ShouldEqual, "bar")
			So(vars[1], ShouldEqual, 10)
		} else if query == "INSERT INTO people (baz, foo) VALUES ($1, $2)" {
			So(vars[0], ShouldEqual, 10)
			So(vars[1], ShouldEqual, "bar")
		} else {
			panic("Invalid query")
		}
	})
}
