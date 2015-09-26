package sqlbuilder

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type Person struct {
	ID        int64  `db:"id"`
	FirstName string `db:"first_name"`
	Birthday  string
}

func TestGetData(t *testing.T) {
	person := &Person{10, "Testy", "1/1/1988"}

	Convey("getData", t, func() {
		Convey("success - struct", func() {
			data, err := getData(person)
			So(err, ShouldBeNil)

			So(len(data.keys), ShouldEqual, 3)
			So(data.data["id"], ShouldEqual, 10)
			So(data.data["first_name"], ShouldEqual, "Testy")
			So(data.data["Birthday"], ShouldEqual, "1/1/1988")

			So(data.keys[0], ShouldEqual, "id")
			So(data.keys[1], ShouldEqual, "first_name")
			So(data.keys[2], ShouldEqual, "Birthday")
		})

		Convey("success - map", func() {
			data, err := getData(map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			})

			So(err, ShouldBeNil)

			So(len(data.keys), ShouldEqual, 3)
			So(data.data["a"], ShouldEqual, 1)
			So(data.data["b"], ShouldEqual, 2)
			So(data.data["c"], ShouldEqual, 3)

			// We explicitly cannot test the order because that's the nature of go maps
		})

		Convey("failure - invalid map", func() {
			data, err := getData(map[int]string{
				1: "a",
				2: "b",
				3: "c",
			})

			So(err, ShouldNotBeNil)
			So(data, ShouldBeNil)
		})
		Convey("failure - invalid argument", func() {
			data, err := getData(10)
			So(err, ShouldNotBeNil)
			So(data, ShouldBeNil)
		})
	})
}
