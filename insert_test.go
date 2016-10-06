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
}

func TestInsert(t *testing.T) {
	person := &Person{10, "Testy", "12345", nil}
	Convey("insert", t, func() {
		query, vars := Insert(person).Into("people").Returning(`"id"`).GetSQL()

		So(len(vars), ShouldEqual, 3)

		// Could be one of three options here...

		options := []string{
			`INSERT INTO people (id, first_name, Birthday) VALUES ($1, $2, $3) RETURNING "id"`,
			`INSERT INTO people (id, Birthday, first_name) VALUES ($1, $2, $3) RETURNING "id"`,
			`INSERT INTO people (first_name, id, Birthday) VALUES ($1, $2, $3) RETURNING "id"`,
			`INSERT INTO people (first_name, Birthday, id) VALUES ($1, $2, $3) RETURNING "id"`,
			`INSERT INTO people (Birthday, first_name, id) VALUES ($1, $2, $3) RETURNING "id"`,
			`INSERT INTO people (Birthday, id, first_name) VALUES ($1, $2, $3) RETURNING "id"`,
		}

		So(options, ShouldContain, query)
		So(vars, ShouldContain, 10)
		So(vars, ShouldContain, "Testy")
		So(vars, ShouldContain, "12345")
	})
}
