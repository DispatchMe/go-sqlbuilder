package sqlbuilder

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// Run UPDATE using a struct. The `db` tags designate alternate column names; otherwise the verbatim struct property will be used. In this case "Age" is the interpreted column name because that property has no `db` tag. Note that Update.Set also supports maps with string keys as its first argument, but because map output is random, it is impossible to test against an expected output (the order of keys is random, but the order of returned variables will always match the order of the keys)
func ExampleUpdate_struct() {
	type Person struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Age       int
	}
	sql, vars := Update("people").Set(&Person{"Testy", "McGee", 25}).Where(Equal{"first_name", "Joe"}).Limit(1).GetSQL()

	fmt.Println(sql, ",", vars)
}

func TestUpdate(t *testing.T) {
	data := map[string]interface{}{
		"foo": "bar",
		"baz": 10,
	}

	Convey("update", t, func() {
		query, vars := Update("people").Set(data).Where(Equal{"id", 100}).Limit(1).GetSQL()

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
