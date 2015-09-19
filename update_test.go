package sqlbuilder

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUpdate(t *testing.T) {
	data := map[string]interface{}{
		"foo": "bar",
		"baz": 10,
	}

	Convey("update", t, func() {
		query, vars := Update("people").Set(data).Where(Expr{"id", "=", 100}).Limit(1).GetSQL()

		So(len(vars), ShouldEqual, 3)
		// Problem here is map output is random, so we need to allow for both possibilities
		if query == "UPDATE people SET foo=$1, baz=$2 WHERE id = $3 LIMIT 1" {
			So(vars[0], ShouldEqual, "bar")
			So(vars[1], ShouldEqual, 10)
		} else if query == "UPDATE people SET baz=$1, foo=$2 WHERE id = $3 LIMIT 1" {
			So(vars[0], ShouldEqual, 10)
			So(vars[1], ShouldEqual, "bar")
		} else {
			panic("Invalid query")
		}

		So(vars[2], ShouldEqual, 100)
	})
}
