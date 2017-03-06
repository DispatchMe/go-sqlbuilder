package sqlbuilder

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// Basic DELETE with a where clause
func ExampleDelete_basic() {
	sql, vars := Delete().From("people").Where(And(
		Equal{"first_name", "Joe"},
		Equal{"last_name", "Blow"},
	)).GetFullSQL()

	fmt.Println(sql, ",", vars)
	// Output: DELETE FROM people WHERE (first_name = $1 AND last_name = $2) , [Joe Blow]
}

// Complex DELETE with JOIN and HAVING clauses. Delete all posts with fewer than 10 comments
func ExampleDelete_complex() {
	sql, vars := Delete().From("posts").LeftJoin("comments", OnColumn("comments.post_id", "posts.id")).GroupBy("posts.id").Having(LessThan{"COUNT(comments)", 10}).GetFullSQL()

	fmt.Println(sql, ",", vars)
	// Output: DELETE FROM posts LEFT JOIN comments ON comments.post_id = posts.id GROUP BY posts.id HAVING COUNT(comments) < $1 , [10]
}

func TestDelete(t *testing.T) {
	Convey("delete", t, func() {
		query, vars := Delete().From("foos").Where(Equal{"id", 10}).Limit(1).GetFullSQL()

		So(len(vars), ShouldEqual, 1)
		So(query, ShouldEqual, "DELETE FROM foos WHERE id = $1 LIMIT 1")
		So(vars[0], ShouldEqual, 10)
	})
}
