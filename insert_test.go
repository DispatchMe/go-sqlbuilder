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
	sql, vars := Insert(&Person{"Testy", "McGee", 25}).Into("people").GetSQL()

	fmt.Println(sql, ",", vars)
	// Output: INSERT INTO people (first_name, last_name, Age) VALUES ($1, $2, $3) , [Testy McGee 25]
}

func TestInsert(t *testing.T) {
	person := &Person{10, "Testy", "12345"}
	Convey("insert", t, func() {
		query, vars := Insert(person).Into("people").Returning(`"id"`).GetSQL()

		So(len(vars), ShouldEqual, 3)
		So(query, ShouldEqual, `INSERT INTO people (id, first_name, Birthday) VALUES ($1, $2, $3) RETURNING "id"`)
		So(vars[0], ShouldEqual, 10)
		So(vars[1], ShouldEqual, "Testy")
		So(vars[2], ShouldEqual, "12345")
	})
}
