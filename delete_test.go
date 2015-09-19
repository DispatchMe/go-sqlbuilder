package sqlbuilder

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDelete(t *testing.T) {
	Convey("delete", t, func() {
		query, vars := Delete().From("foos").Where(Expr{"id", "=", 10}).Limit(1).GetSQL()

		So(len(vars), ShouldEqual, 1)
		So(query, ShouldEqual, "DELETE FROM foos WHERE id = $1 LIMIT 1")
		So(vars[0], ShouldEqual, 10)
	})
}
