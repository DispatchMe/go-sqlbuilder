package sqlbuilder

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// Run INSERT using a struct. The `db` tags designate alternate column names; otherwise the verbatim struct property will be used. In this case "Age" is the interpreted column name because that property has no `db` tag. Note that Insert also supports maps with string keys as its first argument, but because map output is random, it is impossible to test against an expected output (the order of keys is random, but the order of returned variables will always match the order of the keys)
func ExampleInsert_struct() {
	type Person struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Age       int
	}
	sql, vars := Insert(&Person{"Testy", "McGee", 25}).Into("people").GetFullSQL()

	fmt.Println(sql, ",", vars)
}

func TestInsert(t *testing.T) {
	person := &Person{10, "Testy", "12345", nil}
	Convey("insert", t, func() {
		query, vars := Insert(person).Into("people").Returning(`"id"`).GetFullSQL()

		// 4th is the "nil" for "pointer"
		So(len(vars), ShouldEqual, 4)
		So(query, ShouldStartWith, "INSERT INTO people (")
		So(query, ShouldContainSubstring, "id")
		So(query, ShouldContainSubstring, "first_name")
		So(query, ShouldContainSubstring, "Birthday")
		So(query, ShouldContainSubstring, "pointer")
		So(query, ShouldEndWith, `RETURNING "id"`)
		So(vars, ShouldContain, int64(10))
		So(vars, ShouldContain, "Testy")
		So(vars, ShouldContain, "12345")
		So(vars, ShouldContain, (*int)(nil))
	})
}
