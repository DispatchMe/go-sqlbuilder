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
		Convey("success", func() {
			data, err := getData(person)
			So(err, ShouldBeNil)

			So(len(data), ShouldEqual, 3)
			So(data["id"], ShouldEqual, 10)
			So(data["first_name"], ShouldEqual, "Testy")
			So(data["Birthday"], ShouldEqual, "1/1/1988")
		})
		Convey("failure", func() {
			data, err := getData(10)
			So(err, ShouldNotBeNil)
			So(data, ShouldBeNil)
		})
	})
}
